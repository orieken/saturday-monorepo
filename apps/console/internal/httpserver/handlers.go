package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"saturday/console/internal/registry"
	"saturday/console/internal/runner"
	"saturday/console/internal/runs"
)

type Handlers struct {
    reg             *registry.Registry
    runStore        *runs.RunStore
    runner          *runner.CucumberRunner
    defaultExecutor string
}

func NewHandlers(reg *registry.Registry, runStore *runs.RunStore) *Handlers {
    projDir := "." // project root relative to execution dir (if running from root)

    // Read default executor from environment; allowed values: "docker" or "k8s". Default to "docker".
    defExec := strings.ToLower(strings.TrimSpace(os.Getenv("DEFAULT_EXECUTOR")))
    if defExec != "k8s" && defExec != "docker" {
        defExec = "docker"
    }

    return &Handlers{
        reg:             reg,
        runStore:        runStore,
        runner:          runner.NewCucumberRunner(reg, projDir, "reports"),
        defaultExecutor: defExec,
    }
}

// GET /api/frameworks
func (h *Handlers) ListFrameworks(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, []string{"cucumber"})
}

// GET /api/frameworks/{frameworkId}/suites
func (h *Handlers) ListSuites(w http.ResponseWriter, r *http.Request) {
    frameworkId := chi.URLParam(r, "frameworkId")
    suites := h.reg.ListSuites(frameworkId)
    writeJSON(w, suites)
}

// GET /api/frameworks/{frameworkId}/suites/{suiteId}/scenarios
func (h *Handlers) ListScenarios(w http.ResponseWriter, r *http.Request) {
    frameworkId := chi.URLParam(r, "frameworkId")
    suiteId := chi.URLParam(r, "suiteId")

    idx, err := h.reg.GetSuiteIndex(frameworkId, suiteId)
    if err != nil {
        writeError(w, http.StatusNotFound, err.Error())
        return
    }

    writeJSON(w, idx)
}

// GET /api/cucumber/index
// Return all registered Cucumber indexes (for convenience)
func (h *Handlers) GetCucumberIndex(w http.ResponseWriter, r *http.Request) {
    suites := h.reg.ListSuites("cucumber")
    var indexes []*registry.CucumberIndex
    for _, s := range suites {
        if idx, err := h.reg.GetSuiteIndex("cucumber", s); err == nil {
            indexes = append(indexes, idx)
        }
    }
    writeJSON(w, indexes)
}

// POST /api/runs
func (h *Handlers) RunScenario(w http.ResponseWriter, r *http.Request) {
    var req runs.RunRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON")
        return
    }

    // If the request didn't specify an executor, use the configured default
    if strings.TrimSpace(req.Executor) == "" {
        req.Executor = h.defaultExecutor
        fmt.Fprintf(os.Stderr, "no executor specified in request; using default executor=%s for suite=%s scenario=%s\n", req.Executor, req.SuiteId, req.ScenarioId)
    } else {
        fmt.Fprintf(os.Stderr, "executor provided in request: executor=%s for suite=%s scenario=%s\n", req.Executor, req.SuiteId, req.ScenarioId)
    }

    run := h.runStore.StartRun(req)
    // Log run creation for debugging
    fmt.Fprintf(os.Stderr, "run created: %s framework=%s suite=%s scenario=%s executor=%s\n", run.ID, run.Framework, run.SuiteId, run.ScenarioId, run.Executor)

    // Kick off execution asynchronously
    go func() {
        fmt.Fprintf(os.Stderr, "run goroutine starting for: %s\n", run.ID)
        // use a background context so the run continues even after the HTTP request completes
        ctx := context.Background()
        status, reportURL, err := h.runner.Run(ctx, run)
        if err != nil {
            fmt.Fprintf(os.Stderr, "run error: %v\n", err)
        }
        h.runStore.CompleteRun(run.ID, status, reportURL)
    }()

    writeJSON(w, run)
}

// GET /api/runs/{runId}
func (h *Handlers) GetRun(w http.ResponseWriter, r *http.Request) {
    runId := chi.URLParam(r, "runId")
    run, err := h.runStore.GetRun(runId)
    if err != nil {
        writeError(w, http.StatusNotFound, err.Error())
        return
    }
    writeJSON(w, run)
}

// POST /api/cucumber/index
// Accepts parsed feature metadata (index) and registers it.
func (h *Handlers) IngestCucumberIndex(w http.ResponseWriter, r *http.Request) {
    var idx registry.CucumberIndex
    if err := json.NewDecoder(r.Body).Decode(&idx); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON")
        return
    }

    // Persist to disk for debugging / reload
    if err := os.MkdirAll("data", 0o755); err == nil {
        if data, err := json.MarshalIndent(idx, "", "  "); err == nil {
            _ = os.WriteFile(fmt.Sprintf("data/%s-index.json", idx.SuiteId), data, 0o644)
        }
    }

    h.reg.RegisterCucumberIndex(&idx)

    writeJSON(w, map[string]string{
        "status":          "ok",
        "registeredSuite": idx.SuiteId,
    })
}

// GET /api/runs/{runId}/logs
func (h *Handlers) GetRunLogs(w http.ResponseWriter, r *http.Request) {
    runId := chi.URLParam(r, "runId")
    run, err := h.runStore.GetRun(runId)
    if err != nil {
        writeError(w, http.StatusNotFound, err.Error())
        return
    }

    // logs are stored under ./reports/<suiteId>/<runId>/run.log
    logPath := filepath.Join("./reports", run.SuiteId, run.ID, "run.log")
    data, err := os.ReadFile(logPath)
    if err != nil {
        // if file not found, return empty lines
        if os.IsNotExist(err) {
            writeJSON(w, map[string]interface{}{"lines": []string{}, "from": 0})
            return
        }
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }

    content := string(data)
    lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")

    // handle fromLine param
    fromLine := 0
    if q := r.URL.Query().Get("fromLine"); q != "" {
        if v, err := strconv.Atoi(q); err == nil && v >= 0 {
            fromLine = v
        }
    }

    if fromLine >= len(lines) {
        writeJSON(w, map[string]interface{}{"lines": []string{}, "from": len(lines)})
        return
    }

    writeJSON(w, map[string]interface{}{"lines": lines[fromLine:], "from": fromLine})
}

// GET /api/runs/{runId}/stream
func (h *Handlers) StreamRunLogs(w http.ResponseWriter, r *http.Request) {
    runId := chi.URLParam(r, "runId")
    run, err := h.runStore.GetRun(runId)
    if err != nil {
        writeError(w, http.StatusNotFound, err.Error())
        return
    }

    // logs are stored under ./reports/<suiteId>/<runId>/run.log
    logPath := filepath.Join("./reports", run.SuiteId, run.ID, "run.log")

    // Setup SSE headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    flusher, ok := w.(http.Flusher)
    if !ok {
        writeError(w, http.StatusInternalServerError, "streaming unsupported")
        return
    }

    // Start position: 0 to stream full history
    var pos int64 = 0

    // Loop and stream new lines as they appear
    ctx := r.Context()
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // read any new data
            fi, err := os.Stat(logPath)
            if err == nil && fi.Size() > pos {
                f, err := os.Open(logPath)
                if err == nil {
                    _, _ = f.Seek(pos, 0)
                    buf := make([]byte, fi.Size()-pos)
                    n, _ := f.Read(buf)
                    f.Close()
                    pos += int64(n)
                    // split into lines
                    content := string(buf[:n])
                    lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
                    for _, line := range lines {
                        if line == "" {
                            continue
                        }
                        // build JSON message
                        msg := map[string]string{"type": "log", "text": line}
                        js, _ := json.Marshal(msg)
                        // write as SSE data
                        fmt.Fprintf(w, "data: %s\n\n", js)
                        flusher.Flush()
                    }
                }
            }

            // send a status event if run completed
            rcur, _ := h.runStore.GetRun(runId)
            if rcur != nil && rcur.Status != "running" {
                msg := map[string]string{"type": "status", "status": rcur.Status}
                js, _ := json.Marshal(msg)
                fmt.Fprintf(w, "data: %s\n\n", js)
                flusher.Flush()
                return
            }
        }
    }
}

// GET /api/runs/{runId}/report
// Returns the cucumber.json report for the given run
func (h *Handlers) GetRunReport(w http.ResponseWriter, r *http.Request) {
    runId := chi.URLParam(r, "runId")
    run, err := h.runStore.GetRun(runId)
    if err != nil {
        writeError(w, http.StatusNotFound, err.Error())
        return
    }

    // Report is stored under ./reports/<suiteId>/<runId>/cucumber.json
    reportPath := filepath.Join("./reports", run.SuiteId, run.ID, "cucumber.json")
    data, err := os.ReadFile(reportPath)
    if err != nil {
        // if file not found, return empty array
        if os.IsNotExist(err) {
            writeJSON(w, []interface{}{})
            return
        }
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }

    // Parse and re-serialize to ensure valid JSON
    var report interface{}
    if err := json.Unmarshal(data, &report); err != nil {
        writeError(w, http.StatusInternalServerError, "invalid JSON in report file")
        return
    }

    writeJSON(w, report)
}


// GET /api/config
func (h *Handlers) GetConfig(w http.ResponseWriter, r *http.Request) {
    cfg := map[string]string{
        "defaultExecutor": h.defaultExecutor,
    }
    writeJSON(w, cfg)
}

// GET /api/debug/ls
func (h *Handlers) DebugListReports(w http.ResponseWriter, r *http.Request) {
    sub := r.URL.Query().Get("sub") // optional subdir
    root := "./reports"
    if sub != "" {
        root = filepath.Join(root, sub)
    }
    
    entries, err := os.ReadDir(root)
    if err != nil {
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    var names []string
    for _, e := range entries {
        names = append(names, e.Name())
    }
    writeJSON(w, names)
}

// GET /api/debug/cat
func (h *Handlers) DebugReadFile(w http.ResponseWriter, r *http.Request) {
    sub := r.URL.Query().Get("file")
    path := filepath.Join("./reports", sub)
    
    data, err := os.ReadFile(path)
    if err != nil {
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    w.Write(data)
}

// GET /
func (h *Handlers) Root(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Test Runner Service</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background-color: #0f172a; /* Slate 900 */
            color: #f8fafc; /* Slate 50 */
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            height: 100vh;
            margin: 0;
        }
        .container {
            text-align: center;
            padding: 2rem;
            background-color: #1e293b; /* Slate 800 */
            border-radius: 1rem;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
            border: 1px solid #334155; /* Slate 700 */
        }
        img {
            width: 150px;
            height: auto;
            margin-bottom: 1.5rem;
            filter: drop-shadow(0 0 10px rgba(56, 189, 248, 0.3)); /* Sky 400 shadow */
        }
        h1 {
            font-size: 2rem;
            font-weight: 700;
            margin: 0 0 0.5rem 0;
            background: linear-gradient(to right, #60a5fa, #c084fc); /* Blue to Purple gradient */
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        p {
            color: #94a3b8; /* Slate 400 */
            margin: 0;
            font-size: 1.1rem;
        }
        .status-dot {
            display: inline-block;
            width: 10px;
            height: 10px;
            background-color: #4ade80; /* Green 400 */
            border-radius: 50%;
            margin-right: 0.5rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <img src="/logo.png" alt="Test Runner Service Logo">
        <h1>Test Runner Service</h1>
        <p><span class="status-dot"></span>Operational</p>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
