package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riverajo/fitness-app/backend/internal/util" // Our JWT utility
)

// Define a key type for storing the user ID in the context
type ContextKey string

const UserIDKey ContextKey = "user_id"

// AuthMiddleware extracts and verifies the JWT token from the "auth_token" cookie.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Try to read the cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			// If cookie is not present (or bad name), continue without a user ID.
			// This allows unauthenticated queries (like login/register) to proceed.
			next.ServeHTTP(w, r)
			return
		}

		// 2. Validate the token
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			// CRITICAL: Handle missing secret gracefully (or fatal on startup)
			fmt.Println("CRITICAL: JWT_SECRET not set in environment.")
			next.ServeHTTP(w, r)
			return
		}

		tokenString := cookie.Value
		claims := &util.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			// Ensure token signing method is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
			}
			return []byte(jwtSecret), nil
		})

		// 3. Handle validation failures (expired, invalid signature, etc.)
		if err != nil || !token.Valid {
			// Token is invalid/expired. Continue, but user is not authenticated.
			next.ServeHTTP(w, r)
			return
		}

		// 4. Token is valid: Inject UserID into the context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		// 5. Pass the request to the next handler (the GraphQL server)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// ResponseWriterMiddleware attaches the http.ResponseWriter to the context.
func ResponseWriterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use a different context key for the ResponseWriter
		ctx := context.WithValue(r.Context(), "ResponseWriterKey", w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetResponseWriter is a helper to retrieve the ResponseWriter from the context.
func GetResponseWriter(ctx context.Context) http.ResponseWriter {
	return ctx.Value("ResponseWriterKey").(http.ResponseWriter)
}
