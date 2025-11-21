package graph

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/riverajo/fitness-app/backend/internal/middleware"
)

func TestLogoutMutation(t *testing.T) {
	// 1. Setup
	// Logout doesn't use any services, so we can initialize Resolver with nil services
	resolver := &Resolver{}

	// Create a ResponseRecorder to capture the cookie
	w := httptest.NewRecorder()

	// Create a context with the ResponseWriter injected (using the key from middleware)
	ctx := context.WithValue(context.Background(), "ResponseWriterKey", w)

	// 2. Execute
	payload, err := resolver.Mutation().Logout(ctx)

	// 3. Verify
	if err != nil {
		t.Fatalf("Logout returned error: %v", err)
	}

	if !payload.Success {
		t.Error("Expected success to be true")
	}

	if payload.Message != "Logged out successfully." {
		t.Errorf("Unexpected message: %s", payload.Message)
	}

	// Check the cookie
	result := w.Result()
	cookies := result.Cookies()

	found := false
	for _, cookie := range cookies {
		if cookie.Name == middleware.AuthCookieName {
			found = true
			if cookie.MaxAge != -1 {
				t.Errorf("Expected MaxAge to be -1, got %d", cookie.MaxAge)
			}
			if cookie.Value != "" {
				t.Errorf("Expected cookie value to be empty, got %s", cookie.Value)
			}
			break
		}
	}

	if !found {
		t.Error("auth_token cookie was not set")
	}
}
