package middleware

import (
	"context"
	"fmt"
	"net/http"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

// GenerateRefreshToken creates a new refresh token for the user.
// It has a longer expiration time (e.g., 7 days) and may include a specific type claim.
func GenerateRefreshToken(user *model.User, jwtSecret string) (string, error) {
	if jwtSecret == "" {
		return "", fmt.Errorf("jwtSecret is required")
	}

	// Token expiration (7 days)
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// Create the Claims
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// We can use the Subject or a custom claim to distinguish, but for now relying on expiration
			// and context of usage is often enough in simple setups.
			// Adding a custom "type" claim would be better practice if we modify Claims struct.
		},
	}

	// Create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", fmt.Errorf("error signing refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateRefreshToken parses and validates the refresh token.
// It returns the claims if valid, or an error otherwise.
func ValidateRefreshToken(tokenString, jwtSecret string) (*Claims, error) {
	if jwtSecret == "" {
		return nil, fmt.Errorf("jwtSecret is required")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired refresh token: %v", err)
	}

	return claims, nil
}
