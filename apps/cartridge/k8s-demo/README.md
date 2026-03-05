K8s Demo: Running a test job and retrieving logs/reports

This demo shows a minimal Kubernetes Job that simulates a test run, writes log lines to stdout, and writes a simple HTML report to /reports/index.html inside the pod. The demo does not require any special image beyond Alpine and `kubectl`/`kind`/`minikube` to interact with the cluster.

Files
- demo-job.yaml — Job manifest that runs a short script which emits logs and writes a small HTML report.

Quick steps (assumes you have a k8s cluster available via kubectl; `kind` or `minikube` works):

1) Apply the Job

```bash
kubectl apply -f demo-job.yaml
```

2) Watch the Job and Pod

```bash
kubectl get jobs
kubectl get pods --selector=job-name=test-runner-demo-job
```

3) Stream logs from the Pod (replace POD with the pod name reported above)

```bash
kubectl logs -f <POD>
```

You should see lines like:

```
Starting demo run: demo-1
[INFO] Step 1 - running...
[INFO] Step 2 - running...
...
Report written to /reports/index.html
Demo run complete
```

4) Fetch the report from the pod (copy it out)

```bash
kubectl cp <POD>:/reports/index.html ./demo-report-demo-1.html
```

Open `demo-report-demo-1.html` in your browser to view the simple demo report.

Notes and next steps
- The job writes the report to the pod's filesystem. In a real setup you'd have the job upload the report to object storage (S3/MinIO) or a service endpoint so the server can expose `reportUrl`.
- To integrate this into the test-runner-service, implement a Kubernetes executor that creates similar Jobs and streams Pod logs to the server (either by polling the API or streaming via the `pods/log` endpoint). The server should collect the report artifact (via object store upload or `kubectl cp` from a controller) and update the run record.
- For local testing with `kind`, ensure you have `kubectl` configured to talk to the `kind` cluster and that `kubectl cp` works in your environment.

