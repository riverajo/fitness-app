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

		// 1. Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No header, continue unauthenticated
			next.ServeHTTP(w, r)
			return
		}

		// 2. Check format "Bearer <token>"
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			// Invalid format, just continue unauthenticated (or could error if strictly required)
			// For now, treat as no token.
			next.ServeHTTP(w, r)
			return
		}

		tokenString := authHeader[7:]

		// 3. Validate the token
		if jwtSecret == "" {
			fmt.Println("CRITICAL: JWT_SECRET not set in environment.")
			next.ServeHTTP(w, r)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
			}
			return []byte(jwtSecret), nil
		})

		// 4. Handle validation failures
		if err != nil || !token.Valid {
			next.ServeHTTP(w, r)
			return
		}

		// 5. Token is valid: Inject UserID into the context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

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
