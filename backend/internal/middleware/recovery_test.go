package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoveryMiddleware(t *testing.T) {
	// Capture logs
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	slog.SetDefault(logger)

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	handler := RecoveryMiddleware(panicHandler)

	req := httptest.NewRequest("GET", "/panic", nil)
	rr := httptest.NewRecorder()

	// This should not panic
	handler.ServeHTTP(rr, req)

	// Check response code
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusInternalServerError)
	}

	// Check logs for panic recovery
	logOutput := buf.String()
	if !strings.Contains(logOutput, "PANIC RECOVERED") {
		t.Errorf("expected log to contain 'PANIC RECOVERED', got %s", logOutput)
	}
	if !strings.Contains(logOutput, "something went wrong") {
		t.Errorf("expected log to contain panic message, got %s", logOutput)
	}
}
