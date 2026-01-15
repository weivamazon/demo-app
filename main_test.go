package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response Response
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("handler returned unexpected status: got %v want %v", response.Status, "healthy")
	}
}

func TestVersionHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(versionHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var info VersionInfo
	if err := json.Unmarshal(rr.Body.Bytes(), &info); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if info.Version == "" {
		t.Error("version should not be empty")
	}
}

func TestHelloHandler(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		expectedName string
	}{
		{"default", "", "World"},
		{"with name", "?name=Test", "Test"},
		{"with chinese", "?name=测试", "测试"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/hello"+tt.query, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(helloHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			var response HelloResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("failed to parse response: %v", err)
			}

			// Check that the message contains the expected name
			if !strings.Contains(response.Message, tt.expectedName) {
				t.Errorf("handler returned unexpected message: got %v, expected to contain %v", response.Message, tt.expectedName)
			}

			// Check that it's a greeting message
			if !strings.HasPrefix(response.Message, "Hello,") {
				t.Errorf("handler returned unexpected message format: got %v", response.Message)
			}
		})
	}
}

func TestFeatureHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/feature", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(featureHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response FeatureResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	// Feature should not be empty
	if response.Feature == "" {
		t.Error("feature should not be empty")
	}

	if response.Version != Version {
		t.Errorf("handler returned unexpected version: got %v want %v", response.Version, Version)
	}
}

func TestStatusHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(statusHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response StatusResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if response.Status != "running" {
		t.Errorf("handler returned unexpected status: got %v want %v", response.Status, "running")
	}
}

func TestMetricsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(metricsHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response MetricsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if response.GoRoutines <= 0 {
		t.Error("goRoutines should be positive")
	}
}

func TestEchoHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/echo?message=test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(echoHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response EchoResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if response.Echo != "test" {
		t.Errorf("handler returned unexpected echo: got %v want %v", response.Echo, "test")
	}

	if response.Method != "GET" {
		t.Errorf("handler returned unexpected method: got %v want %v", response.Method, "GET")
	}
}

func TestInfoHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(infoHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response InfoResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if response.AppName != "Demo App" {
		t.Errorf("handler returned unexpected appName: got %v want %v", response.AppName, "Demo App")
	}

	if response.Version != Version {
		t.Errorf("handler returned unexpected version: got %v want %v", response.Version, Version)
	}
}

func TestTimeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/time", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(timeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response TimeResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if response.UnixTime <= 0 {
		t.Error("unixTime should be positive")
	}

	if response.ServerTime == "" {
		t.Error("serverTime should not be empty")
	}
}

func TestRandomHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/random", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(randomHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response RandomResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response: %v", err)
	}

	if response.UUID == "" {
		t.Error("uuid should not be empty")
	}

	if response.Color == "" {
		t.Error("color should not be empty")
	}

	if len(response.Dice) != 3 {
		t.Errorf("dice should have 3 elements, got %d", len(response.Dice))
	}
}
