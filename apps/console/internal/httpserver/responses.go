package httpserver

import (
    "encoding/json"
    "net/http"
)

func writeJSON(w http.ResponseWriter, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
    w.WriteHeader(status)
    writeJSON(w, map[string]string{"error": message})
}
