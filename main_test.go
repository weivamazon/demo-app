package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		name     string
		query    string
		expected string
	}{
		{"default", "", "Hello, World! 👋"},
		{"with name", "?name=Test", "Hello, Test! 👋"},
		{"with chinese", "?name=测试", "Hello, 测试! 👋"},
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

			if response.Message != tt.expected {
				t.Errorf("handler returned unexpected message: got %v want %v", response.Message, tt.expected)
			}
		})
	}
}
