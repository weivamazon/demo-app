package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	Version   = "2.4.0-dev"
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

type MetricsResponse struct {
	RequestCount int64   `json:"requestCount"`
	MemoryUsage  string  `json:"memoryUsage"`
	GoRoutines   int     `json:"goRoutines"`
	Uptime       string  `json:"uptime"`
	Timestamp    string  `json:"timestamp"`
}

type EchoResponse struct {
	Echo      string            `json:"echo"`
	Headers   map[string]string `json:"headers"`
	Method    string            `json:"method"`
	Path      string            `json:"path"`
	Timestamp string            `json:"timestamp"`
}

type InfoResponse struct {
	AppName     string `json:"appName"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	GoVersion   string `json:"goVersion"`
	OS          string `json:"os"`
	Arch        string `json:"arch"`
	Timestamp   string `json:"timestamp"`
}

type TimeResponse struct {
	ServerTime   string `json:"serverTime"`
	Timezone     string `json:"timezone"`
	UnixTime     int64  `json:"unixTime"`
	DayOfWeek    string `json:"dayOfWeek"`
	WeekOfYear   int    `json:"weekOfYear"`
	IsWeekend    bool   `json:"isWeekend"`
	Timestamp    string `json:"timestamp"`
}

type RandomResponse struct {
	Number      int      `json:"number"`
	UUID        string   `json:"uuid"`
	Color       string   `json:"color"`
	Quote       string   `json:"quote"`
	LuckyNumber int      `json:"luckyNumber"`
	Dice        []int    `json:"dice"`
	Timestamp   string   `json:"timestamp"`
}

var startTime = time.Now()
var requestCount int64

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
	http.HandleFunc("/api/metrics", metricsHandler)
	http.HandleFunc("/api/echo", echoHandler)
	http.HandleFunc("/api/info", infoHandler)
	http.HandleFunc("/api/time", timeHandler)
	http.HandleFunc("/api/random", randomHandler)
	http.HandleFunc("/", rootHandler)

	log.Printf("Demo App v%s starting on port %s", Version, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	atomic.AddInt64(&requestCount, 1)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <title>Demo App v2.1</title>
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
    <h1>🚀 Demo App <span class="version-badge">v2.4 开发版</span></h1>
    <p>Version: %s</p>
    <p><strong>🆕 v2.4 新功能：</strong> 添加了 Random API 端点，返回随机数据！</p>
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
    <div class="endpoint">
        <strong>GET</strong> <code>/api/feature</code> - 功能展示
    </div>
    <div class="endpoint new-feature">
        <strong>🆕 GET</strong> <code>/api/metrics</code> - 应用指标 (v2.1新增)
    </div>
    <div class="endpoint new-feature">
        <strong>🆕 GET/POST</strong> <code>/api/echo</code> - 请求回显 (v2.1新增)
    </div>
    <div class="endpoint new-feature">
        <strong>🆕 GET</strong> <code>/api/info</code> - 应用详细信息 (v2.2新增)
    </div>
    <div class="endpoint new-feature">
        <strong>🆕 GET</strong> <code>/api/time</code> - 服务器时间信息 (v2.3新增)
    </div>
    <div class="endpoint new-feature">
        <strong>🆕 GET</strong> <code>/api/random</code> - 随机数据生成 (v2.4新增)
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
	atomic.AddInt64(&requestCount, 1)
	response := FeatureResponse{
		Feature:     "新功能展示",
		Description: "这是 v2.1 开发版的功能端点",
		Version:     Version,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	response := MetricsResponse{
		RequestCount: atomic.LoadInt64(&requestCount),
		MemoryUsage:  fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
		GoRoutines:   runtime.NumGoroutine(),
		Uptime:       time.Since(startTime).Round(time.Second).String(),
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	
	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	
	echo := r.URL.Query().Get("message")
	if echo == "" {
		echo = "Hello from Echo API!"
	}
	
	response := EchoResponse{
		Echo:      echo,
		Headers:   headers,
		Method:    r.Method,
		Path:      r.URL.Path,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	
	response := InfoResponse{
		AppName:     "Demo App",
		Version:     Version,
		Description: "A demo application for CI/CD pipeline testing",
		Author:      "CI/CD Platform Team",
		GoVersion:   runtime.Version(),
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	
	now := time.Now()
	_, week := now.ISOWeek()
	dayOfWeek := now.Weekday()
	isWeekend := dayOfWeek == time.Saturday || dayOfWeek == time.Sunday
	
	response := TimeResponse{
		ServerTime:   now.Format("2006-01-02 15:04:05"),
		Timezone:     now.Location().String(),
		UnixTime:     now.Unix(),
		DayOfWeek:    dayOfWeek.String(),
		WeekOfYear:   week,
		IsWeekend:    isWeekend,
		Timestamp:    now.UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestCount, 1)
	
	colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7", "#DDA0DD", "#98D8C8", "#F7DC6F"}
	quotes := []string{
		"代码是写给人看的，顺便能在机器上运行。",
		"先让它工作，再让它正确，最后让它快。",
		"简单是可靠的先决条件。",
		"过早优化是万恶之源。",
		"好的代码是它自己最好的文档。",
	}
	
	// Generate random dice rolls (3 dice)
	dice := make([]int, 3)
	for i := range dice {
		dice[i] = rand.Intn(6) + 1
	}
	
	response := RandomResponse{
		Number:      rand.Intn(1000),
		UUID:        fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", rand.Int63(), rand.Int31()&0xffff, rand.Int31()&0xffff, rand.Int31()&0xffff, rand.Int63()),
		Color:       colors[rand.Intn(len(colors))],
		Quote:       quotes[rand.Intn(len(quotes))],
		LuckyNumber: rand.Intn(100) + 1,
		Dice:        dice,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
