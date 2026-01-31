package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/riverajo/fitness-app/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type TokenService struct {
	repo repository.RefreshTokenRepository
}

func NewTokenService(repo repository.RefreshTokenRepository) *TokenService {
	return &TokenService{repo: repo}
}

// GenerateRefreshToken creates a random token and its hash.
func (s *TokenService) GenerateRefreshToken() (secret string, hash string, err error) {
	// 32 bytes of random data
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("failed to generate random token: %w", err)
	}
	secret = base64.URLEncoding.EncodeToString(b)

	// Hash the token
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash token: %w", err)
	}

	return secret, string(hashBytes), nil
}

// CreateCompositeRefreshToken generates a new token, saves it, and returns "ID.Secret".
func (s *TokenService) CreateCompositeRefreshToken(ctx context.Context, userID string) (string, error) {
	secret, hash, err := s.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	refreshToken := &model.RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, refreshToken); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", refreshToken.ID, secret), nil
}

// ValidateRotate validates the old refresh token, revokes it (rotation), and issues a new one.
// Returns: newCompositeToken, userID, error
func (s *TokenService) ValidateRotate(ctx context.Context, compositeToken string) (string, string, error) {
	// 1. Split "ID.Secret"
	parts := strings.Split(compositeToken, ".")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid token format")
	}
	tokenID, secret := parts[0], parts[1]

	// 2. Find in DB
	storedToken, err := s.repo.FindByID(ctx, tokenID)
	if err != nil {
		return "", "", fmt.Errorf("database error: %w", err)
	}
	if storedToken == nil {
		return "", "", fmt.Errorf("invalid refresh token (not found)")
	}

	// 3. Check Expiration
	if time.Now().After(storedToken.ExpiresAt) {
		// Cleanup (fire and forget or do it now)
		_ = s.repo.Revoke(ctx, tokenID)
		return "", "", fmt.Errorf("refresh token expired")
	}

	// 4. Validate Secret against Hash
	if err := bcrypt.CompareHashAndPassword([]byte(storedToken.TokenHash), []byte(secret)); err != nil {
		// Potential reuse detection could go here (if we kept history)
		return "", "", fmt.Errorf("invalid token signature")
	}

	// 5. Revoke Old Token (Rotation)
	if err := s.repo.Revoke(ctx, tokenID); err != nil {
		return "", "", fmt.Errorf("failed to revoke used token: %w", err)
	}

	// 6. Issue New Token
	newToken, err := s.CreateCompositeRefreshToken(ctx, storedToken.UserID)
	if err != nil {
		return "", "", fmt.Errorf("failed to rotate token: %w", err)
	}

	return newToken, storedToken.UserID, nil
}

// Revoke allows manual revocation (e.g. logout).
func (s *TokenService) Revoke(ctx context.Context, compositeToken string) error {
	parts := strings.Split(compositeToken, ".")
	if len(parts) != 2 {
		return nil // checking valid format not strictly necessary for simple logout, but good practice
	}
	tokenID := parts[0]
	return s.repo.Revoke(ctx, tokenID)
}
