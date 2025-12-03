package middleware

import (
	"context"
	"fmt"
	"net/http"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riverajo/fitness-app/backend/internal/model"
)

// Define a key type for storing the user ID in the context
type ContextKey string

const (
	UserIDKey         ContextKey = "user_id"
	AuthCookieName    string     = "auth_token"
	ResponseWriterKey ContextKey = "ResponseWriterKey"
)

// AuthMiddleware extracts and verifies the JWT token from the "auth_token" cookie.
func AuthMiddleware(next http.Handler, jwtSecret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Try to read the cookie
		cookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			// If cookie is not present (or bad name), continue without a user ID.
			// This allows unauthenticated queries (like login/register) to proceed.
			next.ServeHTTP(w, r)
			return
		}

		// 2. Validate the token
		if jwtSecret == "" {
			// CRITICAL: Handle missing secret gracefully (or fatal on startup)
			fmt.Println("CRITICAL: JWT_SECRET not set in environment.")
			next.ServeHTTP(w, r)
			return
		}

		tokenString := cookie.Value
		claims := &Claims{}

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
		ctx := context.WithValue(r.Context(), ResponseWriterKey, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetResponseWriter is a helper to retrieve the ResponseWriter from the context.
func GetResponseWriter(ctx context.Context) http.ResponseWriter {
	val := ctx.Value(ResponseWriterKey)
	if val == nil {
		return nil
	}
	return val.(http.ResponseWriter)
}

// Define the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new signed JWT for the user.
func GenerateJWT(user *model.User, jwtSecret string) (string, error) {
	if jwtSecret == "" {
		return "", fmt.Errorf("jwtSecret is required")
	}

	// Token expiration (e.g., 24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the Claims
	claims := &Claims{
		UserID: user.ID, // Use the MongoDB ObjectID as the user ID in the token
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}
