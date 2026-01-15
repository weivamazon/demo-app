package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	Version   = "2.5.0-dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
	tracer    trace.Tracer
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
	TraceID   string `json:"traceId,omitempty"`
}

type StatusResponse struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Uptime      string `json:"uptime"`
	Timestamp   string `json:"timestamp"`
	TraceID     string `json:"traceId,omitempty"`
}

type FeatureResponse struct {
	Feature     string `json:"feature"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Timestamp   string `json:"timestamp"`
}

type MetricsResponse struct {
	RequestCount int64  `json:"requestCount"`
	MemoryUsage  string `json:"memoryUsage"`
	GoRoutines   int    `json:"goRoutines"`
	Uptime       string `json:"uptime"`
	Timestamp    string `json:"timestamp"`
}

type EchoResponse struct {
	Echo      string            `json:"echo"`
	Headers   map[string]string `json:"headers"`
	Method    string            `json:"method"`
	Path      string            `json:"path"`
	Timestamp string            `json:"timestamp"`
	TraceID   string            `json:"traceId,omitempty"`
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
	ServerTime string `json:"serverTime"`
	Timezone   string `json:"timezone"`
	UnixTime   int64  `json:"unixTime"`
	DayOfWeek  string `json:"dayOfWeek"`
	WeekOfYear int    `json:"weekOfYear"`
	IsWeekend  bool   `json:"isWeekend"`
	Timestamp  string `json:"timestamp"`
}

type RandomResponse struct {
	Number      int    `json:"number"`
	UUID        string `json:"uuid"`
	Color       string `json:"color"`
	Quote       string `json:"quote"`
	LuckyNumber int    `json:"luckyNumber"`
	Dice        []int  `json:"dice"`
	Timestamp   string `json:"timestamp"`
	TraceID     string `json:"traceId,omitempty"`
}

var startTime = time.Now()
var requestCount int64


// initTracer initializes OpenTelemetry tracer
func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	// Get OTLP endpoint from environment, default to Jaeger
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "jaeger:4318" // Default to Jaeger in Docker network
	}

	// Create OTLP HTTP exporter
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Get service name from environment
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "demo-app"
	}

	// Create resource with service information
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(Version),
			attribute.String("environment", os.Getenv("APP_ENV")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global TracerProvider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}

func main() {
	ctx := context.Background()

	// Initialize OpenTelemetry
	tp, err := initTracer(ctx)
	if err != nil {
		log.Printf("Warning: Failed to initialize tracer: %v", err)
	} else {
		defer func() {
			if err := tp.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down tracer: %v", err)
			}
		}()
		tracer = otel.Tracer("demo-app")
		log.Println("OpenTelemetry tracer initialized successfully")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Create mux with instrumented handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/version", versionHandler)
	mux.HandleFunc("/api/hello", helloHandler)
	mux.HandleFunc("/api/status", statusHandler)
	mux.HandleFunc("/api/feature", featureHandler)
	mux.HandleFunc("/api/metrics", metricsHandler)
	mux.HandleFunc("/api/echo", echoHandler)
	mux.HandleFunc("/api/info", infoHandler)
	mux.HandleFunc("/api/time", timeHandler)
	mux.HandleFunc("/api/random", randomHandler)
	mux.HandleFunc("/", rootHandler)

	// Wrap with OpenTelemetry HTTP instrumentation
	handler := otelhttp.NewHandler(mux, "demo-app",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)

	log.Printf("Demo App v%s starting on port %s", Version, port)
	log.Printf("OpenTelemetry endpoint: %s", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	log.Fatal(http.ListenAndServe(":"+port, handler))
}


func getTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "rootHandler")
		defer span.End()
		span.SetAttributes(attribute.String("handler", "root"))
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	atomic.AddInt64(&requestCount, 1)
	log.Printf("[INFO] Root page accessed, request count: %d, traceId: %s", requestCount, getTraceID(ctx))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
    <title>Demo App v2.5</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; background: #f0f8ff; }
        h1 { color: #2e8b57; }
        .version-badge { background: #2e8b57; color: white; padding: 5px 10px; border-radius: 15px; font-size: 14px; }
        .endpoint { background: #fff; padding: 10px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #2e8b57; }
        code { background: #e0e0e0; padding: 2px 6px; border-radius: 3px; }
        .new-feature { background: #fffacd; border-left-color: #ffa500; }
        .otel-badge { background: #7B68EE; color: white; padding: 2px 8px; border-radius: 10px; font-size: 12px; margin-left: 10px; }
    </style>
</head>
<body>
    <h1>🚀 Demo App <span class="version-badge">v2.5 开发版</span> <span class="otel-badge">OpenTelemetry</span></h1>
    <p>Version: %s</p>
    <p><strong>🆕 v2.5 新功能：</strong> 集成 OpenTelemetry 分布式追踪和结构化日志！</p>
    <h2>Available Endpoints:</h2>
    <div class="endpoint"><strong>GET</strong> <code>/health</code> - Health check</div>
    <div class="endpoint"><strong>GET</strong> <code>/version</code> - Version info</div>
    <div class="endpoint"><strong>GET</strong> <code>/api/hello</code> - Hello World (with tracing)</div>
    <div class="endpoint"><strong>GET</strong> <code>/api/status</code> - Application status</div>
    <div class="endpoint"><strong>GET</strong> <code>/api/feature</code> - 功能展示</div>
    <div class="endpoint"><strong>GET</strong> <code>/api/metrics</code> - 应用指标</div>
    <div class="endpoint"><strong>GET/POST</strong> <code>/api/echo</code> - 请求回显 (with tracing)</div>
    <div class="endpoint"><strong>GET</strong> <code>/api/info</code> - 应用详细信息</div>
    <div class="endpoint"><strong>GET</strong> <code>/api/time</code> - 服务器时间信息</div>
    <div class="endpoint"><strong>GET</strong> <code>/api/random</code> - 随机数据生成 (with tracing)</div>
</body>
</html>`, Version)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Printf("[INFO] Health check, traceId: %s", getTraceID(ctx))
	
	response := Response{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Printf("[INFO] Version info requested, traceId: %s", getTraceID(ctx))
	
	info := VersionInfo{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
	}
	writeJSON(w, http.StatusOK, info)
}


func helloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Create a child span for business logic
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "helloHandler.processGreeting")
		defer span.End()
		
		name := r.URL.Query().Get("name")
		span.SetAttributes(
			attribute.String("greeting.name", name),
			attribute.String("http.method", r.Method),
		)
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	// Simulate some processing time
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

	log.Printf("[INFO] Hello endpoint called, name: %s, traceId: %s", name, getTraceID(ctx))

	response := HelloResponse{
		Message:   fmt.Sprintf("Hello, %s! 👋 (v2.5 with OpenTelemetry)", name),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		TraceID:   getTraceID(ctx),
	}
	writeJSON(w, http.StatusOK, response)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "statusHandler.getStatus")
		defer span.End()
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	uptime := time.Since(startTime)
	log.Printf("[INFO] Status check, env: %s, uptime: %s, traceId: %s", env, uptime.Round(time.Second), getTraceID(ctx))

	response := StatusResponse{
		Status:      "running",
		Environment: env,
		Uptime:      uptime.Round(time.Second).String(),
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		TraceID:     getTraceID(ctx),
	}
	writeJSON(w, http.StatusOK, response)
}

func featureHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	atomic.AddInt64(&requestCount, 1)
	log.Printf("[INFO] Feature endpoint called, traceId: %s", getTraceID(ctx))
	
	response := FeatureResponse{
		Feature:     "OpenTelemetry 集成",
		Description: "这是 v2.5 开发版，支持分布式追踪和结构化日志",
		Version:     Version,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	atomic.AddInt64(&requestCount, 1)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	log.Printf("[INFO] Metrics requested, requestCount: %d, memory: %.2f MB, traceId: %s",
		requestCount, float64(m.Alloc)/1024/1024, getTraceID(ctx))

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
	ctx := r.Context()
	atomic.AddInt64(&requestCount, 1)

	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "echoHandler.processRequest")
		defer span.End()
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		)
	}

	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	echo := r.URL.Query().Get("message")
	if echo == "" {
		echo = "Hello from Echo API with OpenTelemetry!"
	}

	// Simulate database call
	if tracer != nil {
		_, dbSpan := tracer.Start(ctx, "database.query")
		time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)
		dbSpan.SetAttributes(attribute.String("db.system", "postgresql"))
		dbSpan.End()
	}

	log.Printf("[INFO] Echo endpoint, method: %s, message: %s, traceId: %s", r.Method, echo, getTraceID(ctx))

	response := EchoResponse{
		Echo:      echo,
		Headers:   headers,
		Method:    r.Method,
		Path:      r.URL.Path,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		TraceID:   getTraceID(ctx),
	}
	writeJSON(w, http.StatusOK, response)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	atomic.AddInt64(&requestCount, 1)
	log.Printf("[INFO] Info endpoint called, traceId: %s", getTraceID(ctx))

	response := InfoResponse{
		AppName:     "Demo App",
		Version:     Version,
		Description: "A demo application with OpenTelemetry for CI/CD pipeline testing",
		Author:      "CI/CD Platform Team",
		GoVersion:   runtime.Version(),
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	atomic.AddInt64(&requestCount, 1)

	now := time.Now()
	_, week := now.ISOWeek()
	dayOfWeek := now.Weekday()
	isWeekend := dayOfWeek == time.Saturday || dayOfWeek == time.Sunday

	log.Printf("[INFO] Time endpoint called, serverTime: %s, traceId: %s", now.Format(time.RFC3339), getTraceID(ctx))

	response := TimeResponse{
		ServerTime: now.Format("2006-01-02 15:04:05"),
		Timezone:   now.Location().String(),
		UnixTime:   now.Unix(),
		DayOfWeek:  dayOfWeek.String(),
		WeekOfYear: week,
		IsWeekend:  isWeekend,
		Timestamp:  now.UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}


func randomHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	atomic.AddInt64(&requestCount, 1)

	if tracer != nil {
		var span trace.Span
		ctx, span = tracer.Start(ctx, "randomHandler.generateData")
		defer span.End()
		
		// Simulate external API call
		_, apiSpan := tracer.Start(ctx, "external.api.call")
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		apiSpan.SetAttributes(attribute.String("api.name", "random-generator"))
		apiSpan.End()

		// Simulate cache lookup
		_, cacheSpan := tracer.Start(ctx, "cache.lookup")
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		cacheSpan.SetAttributes(
			attribute.String("cache.type", "redis"),
			attribute.Bool("cache.hit", rand.Float32() > 0.5),
		)
		cacheSpan.End()
	}

	colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7", "#DDA0DD", "#98D8C8", "#F7DC6F"}
	quotes := []string{
		"代码是写给人看的，顺便能在机器上运行。",
		"先让它工作，再让它正确，最后让它快。",
		"简单是可靠的先决条件。",
		"过早优化是万恶之源。",
		"好的代码是它自己最好的文档。",
	}

	dice := make([]int, 3)
	for i := range dice {
		dice[i] = rand.Intn(6) + 1
	}

	randomNum := rand.Intn(1000)
	log.Printf("[INFO] Random endpoint called, number: %d, traceId: %s", randomNum, getTraceID(ctx))

	// Simulate occasional errors for testing
	if randomNum > 950 {
		if tracer != nil {
			span := trace.SpanFromContext(ctx)
			span.SetStatus(codes.Error, "Random error for testing")
			span.RecordError(fmt.Errorf("simulated error: random number too high"))
		}
		log.Printf("[ERROR] Simulated error occurred, number: %d, traceId: %s", randomNum, getTraceID(ctx))
	}

	response := RandomResponse{
		Number:      randomNum,
		UUID:        fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", rand.Int63(), rand.Int31()&0xffff, rand.Int31()&0xffff, rand.Int31()&0xffff, rand.Int63()),
		Color:       colors[rand.Intn(len(colors))],
		Quote:       quotes[rand.Intn(len(quotes))],
		LuckyNumber: rand.Intn(100) + 1,
		Dice:        dice,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		TraceID:     getTraceID(ctx),
	}
	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
