package telemetry

import (
	"context"
	"testing"
)

func TestInitOTel(t *testing.T) {
	// Test Case 1: Disabled Alloy (should succeed and return shutdown func)
	t.Run("Disabled Alloy", func(t *testing.T) {
		shutdown, err := InitOTel(context.Background(), "development", false)
		if err != nil {
			t.Fatalf("InitOTel failed: %v", err)
		}
		if shutdown == nil {
			t.Fatal("shutdown function should not be nil")
		}
		// Test shutdown
		if err := shutdown(context.Background()); err != nil {
			t.Errorf("shutdown failed: %v", err)
		}
	})

	// Test Case 2: Enabled Alloy (In CI/Test environment this might try to connect)
	// We mainly want to ensure it doesn't panic or error immediately if config is valid logic-wise.
	// However, since it creates gRPC connections, it might be flaky if we try to actually connect.
	// We'll stick to testing the configuration paths we changed (conditional logic).
	// If we pass "production" and "false", it should behave like dev but maybe with different logging?
	// Actually logic is: if appEnv == "production" && enableAlloy.

	t.Run("Production without Alloy", func(t *testing.T) {
		shutdown, err := InitOTel(context.Background(), "production", false)
		if err != nil {
			t.Fatalf("InitOTel failed: %v", err)
		}
		if shutdown == nil {
			t.Fatal("shutdown function should not be nil")
		}
		defer func() {
			if err := shutdown(context.Background()); err != nil {
				t.Errorf("shutdown failed: %v", err)
			}
		}()
	})
}
