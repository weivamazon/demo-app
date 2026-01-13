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
	Version   = "2.0.0-dev"
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

type StatusResponse struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Uptime      string `json:"uptime"`
	Timestamp   string `json:"timestamp"`
}

type FeatureResponse struct {
	Feature     string `json:"feature"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Timestamp   string `json:"timestamp"`
}

var startTime = time.Now()

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/version", versionHandler)
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/api/feature", featureHandler)
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
    <title>Demo App v2.0</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; background: #f0f8ff; }
        h1 { color: #2e8b57; }
        .version-badge { background: #2e8b57; color: white; padding: 5px 10px; border-radius: 15px; font-size: 14px; }
        .endpoint { background: #fff; padding: 10px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #2e8b57; }
        code { background: #e0e0e0; padding: 2px 6px; border-radius: 3px; }
        .new-feature { background: #fffacd; border-left-color: #ffa500; }
    </style>
</head>
<body>
    <h1>🚀 Demo App <span class="version-badge">v2.0 开发版</span></h1>
    <p>Version: %s</p>
    <p><strong>🆕 新功能：</strong> 添加了 Feature API 端点！</p>
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
    <div class="endpoint">
        <strong>GET</strong> <code>/api/status</code> - Application status with uptime
    </div>
    <div class="endpoint new-feature">
        <strong>🆕 GET</strong> <code>/api/feature</code> - 新功能展示 (v2.0新增)
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
		Message:   fmt.Sprintf("Hello, %s! 👋 (v2.0)", name),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	
	uptime := time.Since(startTime)
	
	response := StatusResponse{
		Status:      "running",
		Environment: env,
		Uptime:      uptime.Round(time.Second).String(),
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func featureHandler(w http.ResponseWriter, r *http.Request) {
	response := FeatureResponse{
		Feature:     "新功能展示",
		Description: "这是 v2.0 开发版新增的功能端点，用于测试版本回滚",
		Version:     Version,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
