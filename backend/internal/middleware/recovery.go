package middleware

import (
	"log/slog"
	"net/http"
	"runtime"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// This is where you get your stack trace
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, false)]

				slog.Error("PANIC RECOVERED",
					"error", err,
					"url", r.URL.Path,
					"stack", string(stack),
				)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
