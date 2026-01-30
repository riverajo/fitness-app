package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riverajo/fitness-app/backend/internal/model"
)

func TestGenerateJWT(t *testing.T) {
	// Setup
	jwtSecret := "testsecret"

	user := &model.User{
		ID:    "user123",
		Email: "test@example.com",
	}

	// Test Success
	tokenString, err := GenerateJWT(user, jwtSecret)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}

	if tokenString == "" {
		t.Error("Expected token string, got empty")
	}

	// Verify token content
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte("testsecret"), nil
	})

	if err != nil || !token.Valid {
		t.Errorf("Generated token is invalid: %v", err)
	}

	if claims.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got '%s'", claims.UserID)
	}

	// Test Missing Secret
	_, err = GenerateJWT(user, "")
	if err == nil {
		t.Error("Expected error when JWT_SECRET is missing, got nil")
	}
}

func TestRefreshTokens(t *testing.T) {
	jwtSecret := "testsecret"
	user := &model.User{ID: "user123"}

	// 1. Test GenerateRefreshToken
	tokenString, err := GenerateRefreshToken(user, jwtSecret)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	if tokenString == "" {
		t.Fatal("Expected refresh token string, got empty")
	}

	// 2. Test ValidateRefreshToken
	claims, err := ValidateRefreshToken(tokenString, jwtSecret)
	if err != nil {
		t.Fatalf("ValidateRefreshToken failed: %v", err)
	}
	if claims.UserID != "user123" {
		t.Errorf("Expected UserID 'user123', got '%s'", claims.UserID)
	}

	// 3. Test Invalid Token
	_, err = ValidateRefreshToken("invalid.token.string", jwtSecret)
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}

	// 4. Verify it's different from Access Token (optional but good sanity check)
	accessToken, _ := GenerateJWT(user, jwtSecret)
	if tokenString == accessToken {
		t.Error("Refresh token and Access token should be different (due to exp/issue time), but are identical")
	}
}

func TestAuthMiddleware(t *testing.T) {
	jwtSecret := "testsecret"

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)
		if userID != nil {
			w.Header().Set("X-User-ID", userID.(string))
		}
	})

	handler := AuthMiddleware(nextHandler, jwtSecret)

	t.Run("No Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Header().Get("X-User-ID") != "" {
			t.Error("Expected no user ID for request without header")
		}
	})

	t.Run("Invalid Token Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Header().Get("X-User-ID") != "" {
			t.Error("Expected no user ID for request with invalid token")
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		user := &model.User{ID: "user123"}
		token, _ := GenerateJWT(user, jwtSecret)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Header().Get("X-User-ID") != "user123" {
			t.Errorf("Expected user ID 'user123', got '%s'", w.Header().Get("X-User-ID"))
		}
	})
}

func TestResponseWriterMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := GetResponseWriter(r.Context())
		if writer == nil {
			t.Error("Expected ResponseWriter in context, got nil")
		}
		if writer != w {
			t.Error("Context ResponseWriter does not match actual ResponseWriter")
		}
	})

	handler := ResponseWriterMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
}
