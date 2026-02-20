package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func FaroProxy(target string, enableAlloy bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If Alloy is disabled, just acknowledge the request and move on.
		// This prevents "403 Forbidden" or "502" errors in the browser console.
		if !enableAlloy || target == "" {
			// Read and discard the body so the connection can be reused
			_, _ = io.Copy(io.Discard, r.Body)
			err := r.Body.Close()
			if err != nil {
				slog.Error("Failed to close request body", "error", err)
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ignored"}`))
			return
		}

		targetURL, err := url.Parse(target)
		if err != nil {
			slog.Error("Failed to parse Faro URL", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Override the Director to prevent it from appending the original path
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			// Explicitly set the path to the target's path, ignoring the incoming path
			req.URL.Path = targetURL.Path
			req.URL.RawPath = targetURL.RawPath
			// Rewrite the Host header
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
			req.Host = targetURL.Host
		}

		// Set a custom error handler for the proxy to log 502s properly
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			slog.Error("Faro proxy upstream unreachable", "target", target, "error", err)
			w.WriteHeader(http.StatusBadGateway)
		}

		proxy.ServeHTTP(w, r)
	})
}
