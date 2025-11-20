package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/riverajo/fitness-app/backend/internal/repository"
)

// UserService handles all user-related business logic and database interaction.
type UserService struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of the UserService.
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// -------------------------------------------------------------------
// CORE BUSINESS LOGIC
// -------------------------------------------------------------------

// It assumes the user object passed in is fully prepared (hashed, timed, etc.).
func (s *UserService) CreateUser(ctx context.Context, input model.User) error {
	// Delegate to repository
	// The repository implementation already checks for existence, but we could move that check here if we wanted "pure" business logic.
	// However, the repository implementation I wrote does the check.
	// Let's rely on the repository for now as it mimics the previous behavior.
	return s.repo.Create(ctx, input)
}

// Placeholder for verification logic (needed for login)
func (s *UserService) VerifyPassword(ctx context.Context, email, password string) (*model.User, error) {
	// 1. Find user by email
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 2. Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials") // Use a generic error message for security
	}

	return user, nil
}

// GetUserByID fetches a user from the database using their ObjectID.
func (s *UserService) GetUserByID(ctx context.Context, idString string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, idString)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, input model.UserUpdateInput) (*model.User, error) {
	// 1. Fetch the existing user (to verify current password)
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user authentication failed: %w", err)
	}

	// 2. Verify the current password (MUST BE PROVIDED)
	if input.CurrentPassword == nil || *input.CurrentPassword == "" {
		return nil, fmt.Errorf("current password is required for update")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(*input.CurrentPassword))
	if err != nil {
		return nil, fmt.Errorf("invalid current password")
	}

	// 3. Update fields
	updated := false

	// Handle New Password Update
	if input.NewPassword != nil && *input.NewPassword != "" {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(*input.NewPassword), bcrypt.DefaultCost)
		if hashErr != nil {
			return nil, fmt.Errorf("failed to hash new password: %w", hashErr)
		}
		user.PasswordHash = string(hashedPassword)
		updated = true
	}

	// Handle Preferred Unit Update
	if input.PreferredUnit != nil && *input.PreferredUnit != "" {
		user.PreferredUnit = *input.PreferredUnit
		updated = true
	}

	// If there are updates to apply
	if updated {
		user.UpdatedAt = time.Now()
		err = s.repo.Update(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	// 4. Return the updated user entity
	return user, nil
}
