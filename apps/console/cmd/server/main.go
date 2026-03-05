package main

import (
    "log"
    "net/http"
    "os"

    httpserver "saturday/console/internal/httpserver"
    "saturday/console/internal/registry"
    "saturday/console/internal/runs"
)

func main() {
    reg := registry.NewRegistry()
    runStore := runs.NewRunStore()

    // Optional: load a pre-generated cucumber_index.json if it exists
    if err := reg.LoadCucumberIndexFromJSON("data/cucumber_index.json"); err != nil {
        log.Printf("warning: could not load initial cucumber_index.json: %v", err)
    }

    h := httpserver.New(reg, runStore)

    port := os.Getenv("PORT")
    if port == "" {
        port = "9001"
    }

    log.Printf("Test runner service listening on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, h))
}
