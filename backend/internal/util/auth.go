package util

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riverajo/fitness-app/backend/internal/model"
)

// Define the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new signed JWT for the user.
func GenerateJWT(user *model.User) (string, error) {
	// 💡 CRUCIAL: Load the secret key from an environment variable (for production security)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
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
