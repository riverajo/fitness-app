package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFaroProxy(t *testing.T) {
	// 1. Mock upstream server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("upstream received"))
	}))
	defer upstream.Close()

	tests := []struct {
		name        string
		target      string
		enableAlloy bool
		wantStatus  int
		wantBody    string
	}{
		{
			name:        "Disabled Alloy",
			target:      upstream.URL,
			enableAlloy: false,
			wantStatus:  http.StatusOK,
			wantBody:    `{"status":"ignored"}`,
		},
		{
			name:        "Empty Target",
			target:      "",
			enableAlloy: true,
			wantStatus:  http.StatusOK, // Should also ignore if target is empty
			wantBody:    `{"status":"ignored"}`,
		},
		{
			name:        "Enabled Alloy - Correct Proxy",
			target:      upstream.URL,
			enableAlloy: true,
			wantStatus:  http.StatusAccepted,
			wantBody:    "upstream received",
		},
		{
			name:        "Enabled Alloy - Invalid Target",
			target:      "://invalid-url", // trigger url.Parse error
			enableAlloy: true,
			wantStatus:  http.StatusInternalServerError,
			wantBody:    "Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := FaroProxy(tt.target, tt.enableAlloy)

			// Simple POST request
			req := httptest.NewRequest("POST", "/collect", strings.NewReader(`{"logs":[]}`))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantStatus)
			}

			body, _ := io.ReadAll(rr.Body)
			if !strings.Contains(string(body), tt.wantBody) {
				t.Errorf("handler returned unexpected body: got %s want %s", string(body), tt.wantBody)
			}
		})
	}
}
