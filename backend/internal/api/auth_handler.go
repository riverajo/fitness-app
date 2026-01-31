package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/riverajo/fitness-app/backend/internal/middleware"
	"github.com/riverajo/fitness-app/backend/internal/service"
)

type AuthHandler struct {
	TokenService *service.TokenService
	UserService  *service.UserService
	JWTSecret    string
	SecureCookie bool
}

func NewAuthHandler(tokenService *service.TokenService, userService *service.UserService, jwtSecret string, secureCookie bool) *AuthHandler {
	return &AuthHandler{
		TokenService: tokenService,
		UserService:  userService,
		JWTSecret:    jwtSecret,
		SecureCookie: secureCookie,
	}
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// 1. Get Refresh Token from Cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Missing refresh token", http.StatusUnauthorized)
		return
	}
	refreshToken := cookie.Value

	// 2. Validate and Rotate Token
	newRefreshToken, userID, err := h.TokenService.ValidateRotate(r.Context(), refreshToken)
	if err != nil {
		// Log error internally if needed
		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/auth/refresh",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
		})
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// 3. Fetch User Details
	user, err := h.UserService.GetUserByID(r.Context(), userID)
	if err != nil {
		slog.Warn("Failed to fetch user during refresh", "user_id", userID, "error", err)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// 4. Issue New Access Token
	token, err := middleware.GenerateJWT(user, h.JWTSecret)
	if err != nil {
		slog.Error("Failed to generate JWT during refresh", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// 5. Set New Refresh Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   h.SecureCookie,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   7 * 24 * 3600,
	})

	// 6. Return JSON Response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"token":   token,
		"success": true,
		"user":    user, // Optional: return user data update
	}); err != nil {
		slog.Error("Failed to encode refresh response", "error", err)
	}
}
