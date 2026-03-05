package httpserver

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"saturday/console/internal/registry"
	"saturday/console/internal/runs"
)

// redirectFixWriter wraps http.ResponseWriter to fix relative redirects from http.FileServer
type redirectFixWriter struct {
	http.ResponseWriter
	request      *http.Request
	originalPath string
}

// WriteHeader intercepts 301/302 redirects and makes them absolute
func (rw *redirectFixWriter) WriteHeader(code int) {
	if code == http.StatusMovedPermanently || code == http.StatusFound {
		if loc := rw.Header().Get("Location"); loc != "" {
			// If it's a relative redirect, make it absolute
			if !strings.HasPrefix(loc, "http") && !strings.HasPrefix(loc, "/") {
				// Convert relative path to absolute by appending to the original path
				basePath := rw.originalPath
				if !strings.HasSuffix(basePath, "/") {
					basePath += "/"
				}
				rw.Header().Set("Location", basePath+loc)
			}
		}
	}
	rw.ResponseWriter.WriteHeader(code)
}

// New builds the HTTP handler with all routes registered.
func New(reg *registry.Registry, runStore *runs.RunStore) http.Handler {
    h := NewHandlers(reg, runStore)

    r := chi.NewRouter()
    r.Use(middleware.Logger)

    // Simple CORS for dev/demo
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
            w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
            if r.Method == http.MethodOptions {
                w.WriteHeader(http.StatusNoContent)
                return
            }
            next.ServeHTTP(w, r)
        })
    })

    // API routes
    r.Get("/api/frameworks", h.ListFrameworks)
    r.Get("/api/frameworks/{frameworkId}/suites", h.ListSuites)
    r.Get("/api/frameworks/{frameworkId}/suites/{suiteId}/scenarios", h.ListScenarios)

    // Provide full index for convenience
    r.Get("/api/cucumber/index", h.GetCucumberIndex)
    
    // Root page with logo
    r.Get("/", h.Root)
    r.Get("/logo.png", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "logo.png")
    })

    r.Post("/api/runs", h.RunScenario)
    r.Get("/api/runs/{runId}", h.GetRun)
    r.Get("/api/runs/{runId}/logs", h.GetRunLogs)
    r.Get("/api/runs/{runId}/stream", h.StreamRunLogs)
    r.Get("/api/runs/{runId}/report", h.GetRunReport)

    // New: expose server config so clients can read default executor, etc.
    r.Get("/api/config", h.GetConfig)

    // Ingest parsed cucumber index
    r.Post("/api/cucumber/index", h.IngestCucumberIndex)

    // Static reports (HTML/JSON/etc.)
    // Serve files directly to avoid http.FileServer redirect issues
    r.Handle("/reports/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Strip the /reports prefix
        path := r.URL.Path
        if !strings.HasPrefix(path, "/reports/") {
            http.NotFound(w, r)
            return
        }
        
        // Convert URL path to filesystem path
        fsPath := filepath.Join("/app/reports", strings.TrimPrefix(path, "/reports"))
        
        // Check if it's a file or directory
        info, err := os.Stat(fsPath)
        if err != nil {
            http.NotFound(w, r)
            return
        }
        
        // If it's a file, serve it directly
        if !info.IsDir() {
            http.ServeFile(w, r, fsPath)
            return
        }
        
        // If it's a directory, use FileServer
        r.URL.Path = strings.TrimPrefix(path, "/reports")
        http.FileServer(http.Dir("/app/reports")).ServeHTTP(w, r)
    }))

    return r
}
