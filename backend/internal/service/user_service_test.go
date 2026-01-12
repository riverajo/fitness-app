package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/riverajo/fitness-app/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		input := model.User{Email: "test@example.com"}
		mockRepo.On("Create", ctx, input).Return(nil).Once()

		err := service.CreateUser(ctx, input)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error from Repository", func(t *testing.T) {
		input := model.User{Email: "test@example.com"}
		mockRepo.On("Create", ctx, input).Return(errors.New("db error")).Once()

		err := service.CreateUser(ctx, input)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestVerifyPassword(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	email := "test@example.com"

	t.Run("Success", func(t *testing.T) {
		user := &model.User{Email: email, PasswordHash: string(hashedPassword)}
		mockRepo.On("FindByEmail", ctx, email).Return(user, nil).Once()

		result, err := service.VerifyPassword(ctx, email, password)
		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo.On("FindByEmail", ctx, email).Return(nil, nil).Once()

		result, err := service.VerifyPassword(ctx, email, password)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		user := &model.User{Email: email, PasswordHash: string(hashedPassword)}
		mockRepo.On("FindByEmail", ctx, email).Return(user, nil).Once()

		result, err := service.VerifyPassword(ctx, email, "wrong")
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()
	id := "user-1"

	t.Run("Success", func(t *testing.T) {
		user := &model.User{ID: id}
		mockRepo.On("FindByID", ctx, id).Return(user, nil).Once()

		result, err := service.GetUserByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, id).Return(nil, nil).Once()

		result, err := service.GetUserByID(ctx, id)
		assert.Error(t, err) // Expect error when user is nil (as per implementation)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()
	id := "user-1"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	t.Run("Success New Password", func(t *testing.T) {
		user := &model.User{ID: id, PasswordHash: string(hashedPassword), UpdatedAt: time.Now().Add(-1 * time.Hour)}
		newPass := "newpass"
		input := model.UserUpdateInput{
			CurrentPassword: &password,
			NewPassword:     &newPass,
		}

		mockRepo.On("FindByID", ctx, id).Return(user, nil).Once()

		mockRepo.On("Update", ctx, mock.MatchedBy(func(u *model.User) bool {
			// Verify key fields are present
			return u.ID == id && u.PasswordHash != string(hashedPassword) && !u.UpdatedAt.IsZero()
		})).Return(nil).Once()

		result, err := service.UpdateUser(ctx, id, input)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Verify password changed
		err = bcrypt.CompareHashAndPassword([]byte(result.PasswordHash), []byte(newPass))
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Success PreferredUnit Update", func(t *testing.T) {
		user := &model.User{
			ID:            id,
			PasswordHash:  string(hashedPassword),
			UpdatedAt:     time.Now().Add(-1 * time.Hour),
			PreferredUnit: model.WeightUnitKilograms, // Use constants
		}
		newUnit := model.WeightUnitPounds // Use constants
		input := model.UserUpdateInput{
			CurrentPassword: &password,
			PreferredUnit:   &newUnit,
		}

		mockRepo.On("FindByID", ctx, id).Return(user, nil).Once()

		mockRepo.On("Update", ctx, mock.MatchedBy(func(u *model.User) bool {
			return u.ID == id && u.PreferredUnit == newUnit && !u.UpdatedAt.IsZero()
		})).Return(nil).Once()

		result, err := service.UpdateUser(ctx, id, input)
		assert.NoError(t, err)
		assert.Equal(t, newUnit, result.PreferredUnit)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Missing Current Password", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, id).Return(&model.User{ID: id}, nil).Once()

		input := model.UserUpdateInput{CurrentPassword: nil}
		result, err := service.UpdateUser(ctx, id, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "current password is required")
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("Wrong Current Password", func(t *testing.T) {
		user := &model.User{ID: id, PasswordHash: string(hashedPassword)}
		wrongPass := "wrong"
		input := model.UserUpdateInput{CurrentPassword: &wrongPass}

		mockRepo.On("FindByID", ctx, id).Return(user, nil).Once()

		result, err := service.UpdateUser(ctx, id, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid current password")
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "Update")
	})
}
