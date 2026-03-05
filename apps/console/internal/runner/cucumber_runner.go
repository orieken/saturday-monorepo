package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"saturday/console/internal/registry"
	"saturday/console/internal/runs"
)

// CucumberRunner knows how to execute a single CucumberJS + Playwright scenario.
type CucumberRunner struct {
    reg        *registry.Registry
    reportsDir string
    projectDir string
}

func NewCucumberRunner(reg *registry.Registry, projectDir, reportsDir string) *CucumberRunner {
    return &CucumberRunner{
        reg:        reg,
        reportsDir: reportsDir,
        projectDir: projectDir,
    }
}

// helper: build k8s rest config. Try in-cluster and fall back to KUBECONFIG env.
func buildKubeConfig() (*rest.Config, error) {
    cfg, err := rest.InClusterConfig()
    if err == nil {
        return cfg, nil
    }
    // fall back to KUBECONFIG
    if k := os.Getenv("KUBECONFIG"); k != "" {
        cfg, err = clientcmd.BuildConfigFromFlags("", k)
        if err == nil {
            return cfg, nil
        }
    }
    // lastly try default kubeconfig location
    home := os.Getenv("HOME")
    if home != "" {
        k := filepath.Join(home, ".kube", "config")
        if _, err := os.Stat(k); err == nil {
            cfg, err = clientcmd.BuildConfigFromFlags("", k)
            if err == nil {
                return cfg, nil
            }
        }
    }
    return nil, fmt.Errorf("could not build Kubernetes config: %w", err)
}

// Run executes the scenario for the given run and returns final status and report URL.
func (r *CucumberRunner) Run(ctx context.Context, run *runs.Run) (string, string, error) {
    idx, err := r.reg.GetSuiteIndex(run.Framework, run.SuiteId)
    if err != nil {
        return "failed", "", fmt.Errorf("cannot resolve suite: %w", err)
    }

    var featureFile string
    var line int
    var scenarioName string

    found := false
    for _, f := range idx.Features {
        // Check if the ID matches the feature itself (to run the whole feature)
        if f.Id == run.ScenarioId || f.File == run.ScenarioId {
            featureFile = f.File
            line = 0 // 0 indicates whole feature
            scenarioName = f.Name
            found = true
            break
        }

        // Check scenarios
        for _, s := range f.Scenarios {
            if s.Id == run.ScenarioId {
                featureFile = f.File
                line = s.Line
                scenarioName = s.Name
                found = true
                break
            }
        }
        if found {
            break
        }
    }

    if !found {
        return "failed", "", fmt.Errorf("scenario or feature %s not found in suite %s", run.ScenarioId, run.SuiteId)
    }

    // Reports folder on the host (where the service is running)
    hostReportsBase, err := filepath.Abs(r.reportsDir)
    if err != nil {
        return "failed", "", fmt.Errorf("failed to resolve absolute reports path: %w", err)
    }
    hostRunDir := filepath.Join(hostReportsBase, run.SuiteId, run.ID)
    if err := os.MkdirAll(hostRunDir, 0o755); err != nil {
        return "failed", "", fmt.Errorf("cannot create reports dir on host: %w", err)
    }

    // Inside the cucumber-project container, we'll write reports to /app/reports
    containerReportsBase := "/app/reports"

    // Prepare report paths inside container
    jsonReport := filepath.Join(containerReportsBase, run.SuiteId, run.ID, "cucumber.json")
    htmlReport := filepath.Join(containerReportsBase, run.SuiteId, run.ID, "index.html")

    // Feature location argument
    featurePathInContainer := filepath.Join("/app/features", featureFile)
    locationArg := featurePathInContainer
    if line > 0 {
        locationArg = fmt.Sprintf("%s:%d", featurePathInContainer, line)
    }

    // Diagnostic: write the docker command to run.log even if docker fails
    logPath := filepath.Join(hostRunDir, "run.log")
    // Note: command string will be written later when args are fully constructed

    // Only require docker CLI when not using k8s executor
    if run.Executor != "k8s" {
        if _, err := exec.LookPath("docker"); err != nil {
            msg := fmt.Sprintf("docker CLI not found in test-runner-service container: %v\n", err)
            _ = os.WriteFile(logPath, []byte(msg), 0o644)
            return "failed", "", fmt.Errorf("%s", msg)
        }
    }

    // Kubernetes executor using client-go
    if run.Executor == "k8s" {
        ns := "test-runner"
        jobName := fmt.Sprintf("test-runner-job-%s", run.ID[:8])

        // Build kube client
        cfg, err := buildKubeConfig()
        if err != nil {
            _ = os.WriteFile(logPath, []byte(fmt.Sprintf("kube config error: %v\n", err)), 0o644)
            return "failed", "", err
        }
        clientset, err := kubernetes.NewForConfig(cfg)
        if err != nil {
            _ = os.WriteFile(logPath, []byte(fmt.Sprintf("k8s client error: %v\n", err)), 0o644)
            return "failed", "", err
        }

        // Build job object
        job := &batchv1.Job{
            ObjectMeta: metav1.ObjectMeta{
                Name:      jobName,
                Namespace: ns,
                Labels: map[string]string{
                    "run-id": run.ID,
                },
            },
            Spec: batchv1.JobSpec{
                BackoffLimit: func(i int32) *int32 { return &i }(0),
                Template: corev1.PodTemplateSpec{
                    Spec: corev1.PodSpec{
                        RestartPolicy: corev1.RestartPolicyNever,
                        Containers: []corev1.Container{
                            {
                                Name:  "runner",
                                Image: func() string {
                                    if img := os.Getenv("CUCUMBER_IMAGE"); img != "" {
                                        return img
                                    }
                                    return "cucumber-project:local"
                                }(),
                                Command: func() []string {
                                    args := []string{
                                        "npm", "test", "--",
                                    }
                                    if line > 0 {
                                        args = append(args, "--name", scenarioName)
                                    }
                                    args = append(args,
                                        "--format", fmt.Sprintf("json:%s", jsonReport),
                                        "--format", fmt.Sprintf("html:%s", htmlReport),
                                        locationArg,
                                    )
                                    return args
                                }(),
                                Env: []corev1.EnvVar{
                                    {Name: "BASE_URL", Value: "http://web-app:8000"},
                                    {Name: "BROWSER", Value: "chromium"},
                                    {Name: "HEADLESS", Value: "true"},
                                    {Name: "CUCUMBER_HTML_REPORT", Value: htmlReport},
                                    {Name: "CUCUMBER_JSON_REPORT", Value: jsonReport},
                                    {Name: "ENABLE_OTEL", Value: "false"},
                                },
                                VolumeMounts: []corev1.VolumeMount{
                                    {
                                        Name:      "reports",
                                        MountPath: "/app/reports",
                                    },
                                },
                            },
                        },
                        Volumes: []corev1.Volume{
                            {
                                Name: "reports",
                                VolumeSource: corev1.VolumeSource{
                                    PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
                                        ClaimName: "reports-pvc",
                                    },
                                },
                            },
                        },
                    },
                },
            },
        }

        // Create Job
        _, err = clientset.BatchV1().Jobs(ns).Create(ctx, job, metav1.CreateOptions{})
        if err != nil {
            _ = os.WriteFile(logPath, []byte(fmt.Sprintf("failed to create job: %v\n", err)), 0o644)
            return "failed", "", err
        }

        // Poll for pod
        var podName string
        for i := 0; i < 60; i++ {
            podList, _ := clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("job-name=%s", jobName)})
            if podList != nil && len(podList.Items) > 0 {
                pod := podList.Items[0]
                podName = pod.Name
                break
            }
            time.Sleep(1 * time.Second)
        }
        if podName == "" {
            _ = os.WriteFile(logPath, []byte("failed to find pod for job\n"), 0o644)
            return "failed", "", fmt.Errorf("failed to find pod for job %s", jobName)
        }

        // Create the log file immediately so SSE endpoint can start reading
        _ = os.WriteFile(logPath, []byte(fmt.Sprintf("[runner] Starting test run for pod: %s\n", podName)), 0o644)

        // Stream pod logs to run.log
        go func() {
            // open file for append, creating if not exists
            f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
            if err != nil {
                fmt.Fprintf(os.Stderr, "failed to open log file %s: %v\n", logPath, err)
                return
            }
            defer f.Close()

            // follow logs
            req := clientset.CoreV1().Pods(ns).GetLogs(podName, &corev1.PodLogOptions{Follow: true, Container: "runner"})
            stream, err := req.Stream(ctx)
            if err != nil {
                _, _ = f.WriteString(fmt.Sprintf("failed to stream logs: %v\n", err))
                return
            }
            defer stream.Close()
            _, _ = io.Copy(f, stream)
        }()

        // Wait for job completion (succeeded/failed)
        // Use polling on job status
        var finalStatus string = "failed"
        timeout := time.After(600 * time.Second)
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
    LOOP:
        for {
            select {
            case <-ctx.Done():
                break LOOP
            case <-timeout:
                _ = os.WriteFile(logPath, []byte("timeout waiting for job completion\n"), 0o644)
                break LOOP
            case <-ticker.C:
                j, err := clientset.BatchV1().Jobs(ns).Get(ctx, jobName, metav1.GetOptions{})
                if err != nil {
                    // continue polling
                    continue
                }
                if j.Status.Succeeded > 0 {
                    finalStatus = "passed"
                    break LOOP
                }
                if j.Status.Failed > 0 && j.Status.Active == 0 {
                    finalStatus = "failed"
                    break LOOP
                }
            }
        }

        // cleanup job (best-effort)
        // _ = clientset.BatchV1().Jobs(ns).Delete(context.Background(), jobName, metav1.DeleteOptions{PropagationPolicy: func() *metav1.DeletionPropagation { p := metav1.DeletePropagationBackground; return &p }()})

        reportURL := fmt.Sprintf("/reports/%s/%s/index.html", run.SuiteId, run.ID)
        return finalStatus, reportURL, nil
    }

    // Docker executor path
    // This is intended for LOCAL usage (running the service binary on host), not inside K8s.
    
    // Resolve absolute path to reports dir for bind mounting
    absReportsDir, err := filepath.Abs(r.reportsDir)
    if err != nil {
        return "failed", "", fmt.Errorf("failed to resolve absolute reports path: %w", err)
    }

    // Ensure the specific run dir exists locally
    // hostRunDir is already created at the top of the function, but we verify it here.
    if err := os.MkdirAll(hostRunDir, 0o755); err != nil {
        return "failed", "", fmt.Errorf("cannot create local run dir: %w", err)
    }

    // For local Docker execution, we use host networking so the container can reach 
    // localhost:8000 (web-app) and localhost:8001 (mock-api).
    // We bind-mount the reports directory so we can see the results immediately.
    
    runnerImage := os.Getenv("CUCUMBER_IMAGE")
    if runnerImage == "" {
        runnerImage = "cucumber-project:local"
    }

    args := []string{
        "run", "--rm",
        "--network", "host", // Use host network to access local services
        "-v", fmt.Sprintf("%s:%s", absReportsDir, containerReportsBase),
        "-e", "BASE_URL=http://host.docker.internal:8003", // Access host service from container on Mac/Windows
        "-e", "BROWSER=chromium",
        "-e", "HEADLESS=true",
        runnerImage,
        "npm", "test", "--",
    }
    if line > 0 {
        args = append(args, "--name", scenarioName)
    }
    args = append(args,
        "--format", fmt.Sprintf("json:%s", jsonReport),
        "--format", fmt.Sprintf("html:%s", htmlReport),
        locationArg,
    )

    var cmd *exec.Cmd

    cmd = exec.CommandContext(ctx, "docker", args...)
    cmd.Env = append(os.Environ(), fmt.Sprintf("NODE_ENV=%s", "test"))

    // Log the command for debugging
    cmdString := "docker " + strings.Join(args, " ")
    _ = os.WriteFile(logPath, []byte(fmt.Sprintf("[runner] executing: %s\n", cmdString)), 0o644)

    output, err := cmd.CombinedOutput()
    
    // Write output to log file
    f, _ := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0o644)
    if f != nil {
        defer f.Close()
        if len(output) > 0 {
            _, _ = f.WriteString("\n--- docker output ---\n")
            _, _ = f.Write(output)
            _, _ = f.WriteString("\n--- end docker output ---\n")
        }
        if err != nil {
            _, _ = f.WriteString(fmt.Sprintf("docker exec error: %v\n", err))
        }
    }

    status := "passed"
    if err != nil {
        status = "failed"
    }

    reportURL := fmt.Sprintf("/reports/%s/%s/index.html", run.SuiteId, run.ID)
    return status, reportURL, nil
}
