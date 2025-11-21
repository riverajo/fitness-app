package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riverajo/fitness-app/backend/internal/model"
)

func TestGenerateJWT(t *testing.T) {
	// Setup
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	user := &model.User{
		ID:    "user123",
		Email: "test@example.com",
	}

	// Test Success
	tokenString, err := GenerateJWT(user)
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
	os.Unsetenv("JWT_SECRET")
	_, err = GenerateJWT(user)
	if err == nil {
		t.Error("Expected error when JWT_SECRET is missing, got nil")
	}
}

func TestAuthMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDKey)
		if userID != nil {
			w.Header().Set("X-User-ID", userID.(string))
		}
	})

	handler := AuthMiddleware(nextHandler)

	t.Run("No Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Header().Get("X-User-ID") != "" {
			t.Error("Expected no user ID for request without cookie")
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: AuthCookieName, Value: "invalidtoken"})
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Header().Get("X-User-ID") != "" {
			t.Error("Expected no user ID for request with invalid token")
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		user := &model.User{ID: "user123"}
		token, _ := GenerateJWT(user)

		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: AuthCookieName, Value: token})
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
