package spa

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestSPAHandler(t *testing.T) {
	// Mock filesystem
	fs := fstest.MapFS{
		"index.html":       {Data: []byte("<html>index</html>")},
		"assets/app.js":    {Data: []byte("console.log('app')")},
		"assets/style.css": {Data: []byte("body { color: red; }")},
	}

	handler := NewHandler(fs, "index.html")

	tests := []struct {
		name            string
		path            string
		wantStatus      int
		wantBody        string
		wantContentType string
		wantCache       string
	}{
		{
			name:       "Serve existing file",
			path:       "/assets/app.js",
			wantStatus: http.StatusOK,
			wantBody:   "console.log('app')",
		},
		{
			name:            "Serve index.html for root",
			path:            "/",
			wantStatus:      http.StatusOK,
			wantBody:        "<html>index</html>",
			wantContentType: "text/html; charset=utf-8",
			wantCache:       "no-cache, no-store, must-revalidate",
		},
		{
			name:            "Serve index.html for unknown path (client-side routing)",
			path:            "/unknown/route",
			wantStatus:      http.StatusOK,
			wantBody:        "<html>index</html>",
			wantContentType: "text/html; charset=utf-8",
			wantCache:       "no-cache, no-store, must-revalidate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.wantStatus)
			}

			if !strings.Contains(rr.Body.String(), tt.wantBody) {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.wantBody)
			}

			if tt.wantContentType != "" {
				if got := rr.Header().Get("Content-Type"); got != tt.wantContentType {
					t.Errorf("handler returned wrong Content-Type: got %v want %v", got, tt.wantContentType)
				}
			}

			if tt.wantCache != "" {
				if got := rr.Header().Get("Cache-Control"); got != tt.wantCache {
					t.Errorf("handler returned wrong Cache-Control: got %v want %v", got, tt.wantCache)
				}
			}
		})
	}
}

func TestSPAHandler_MissingIndex(t *testing.T) {
	// FS without index.html
	fs := fstest.MapFS{
		"assets/app.js": {Data: []byte("console.log('app')")},
	}

	handler := NewHandler(fs, "index.html")

	req := httptest.NewRequest("GET", "/unknown", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 when index.html is missing, got %d", rr.Code)
	}
}
