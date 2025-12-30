package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	slog.SetDefault(logger)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte("OK"))
	})

	handler := LoggingMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test-path", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Check response
	if rr.Code != http.StatusTeapot {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusTeapot)
	}

	// Check logs
	logOutput := buf.String()
	if !contains(logOutput, "Request completed") {
		t.Errorf("expected log to contain 'Request completed', got %s", logOutput)
	}
	if !contains(logOutput, "method=GET") {
		t.Errorf("expected log to contain method=GET, got %s", logOutput)
	}
	if !contains(logOutput, "path=/test-path") {
		t.Errorf("expected log to contain path=/test-path, got %s", logOutput)
	}
	if !contains(logOutput, "status=418") {
		t.Errorf("expected log to contain status=418, got %s", logOutput)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && bytes.Contains([]byte(s), []byte(substr))
}
