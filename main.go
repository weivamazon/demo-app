package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	Version   = "1.0.0"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

type Response struct {
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
	Timestamp string `json:"timestamp"`
}

type VersionInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
	GitCommit string `json:"gitCommit"`
}

type HelloResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/version", versionHandler)
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/", rootHandler)

	log.Printf("Demo App v%s starting on port %s", Version, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <title>Demo App</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
        code { background: #e0e0e0; padding: 2px 6px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>🚀 Demo App</h1>
    <p>Version: %s</p>
    <h2>Available Endpoints:</h2>
    <div class="endpoint">
        <strong>GET</strong> <code>/health</code> - Health check
    </div>
    <div class="endpoint">
        <strong>GET</strong> <code>/version</code> - Version info
    </div>
    <div class="endpoint">
        <strong>GET</strong> <code>/api/hello</code> - Hello World
    </div>
    <div class="endpoint">
        <strong>GET</strong> <code>/api/hello?name=YourName</code> - Personalized greeting
    </div>
</body>
</html>`, Version)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	info := VersionInfo{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
	}
	writeJSON(w, http.StatusOK, info)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	response := HelloResponse{
		Message:   fmt.Sprintf("Hello, %s! 👋", name),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
