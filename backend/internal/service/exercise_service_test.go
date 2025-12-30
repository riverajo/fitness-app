package service

import (
	"context"
	"errors"
	"testing"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/riverajo/fitness-app/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateExercise(t *testing.T) {
	mockRepo := new(repository.MockExerciseRepository)
	service := NewExerciseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		name := "Push Up"
		desc := "Standard push up"
		userID := "user-123"

		mockRepo.On("Create", ctx, mock.MatchedBy(func(e *model.UniqueExercise) bool {
			return e.Name == name && *e.UserID == userID && *e.Description == desc
		})).Return(nil).Once()

		result, err := service.CreateExercise(ctx, name, &desc, &userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty name", func(t *testing.T) {
		name := "   "
		userID := "user-123"

		result, err := service.CreateExercise(ctx, name, nil, &userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "exercise name cannot be empty", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("repo error", func(t *testing.T) {
		name := "Pull Up"
		userID := "user-123"

		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.UniqueExercise")).Return(errors.New("db error")).Once()

		result, err := service.CreateExercise(ctx, name, nil, &userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "db error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestSearchExercises(t *testing.T) {
	mockRepo := new(repository.MockExerciseRepository)
	service := NewExerciseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		query := "press"
		limit := 10
		offset := 0
		userID := "user-123"
		expected := []*model.UniqueExercise{{Name: "Bench Press"}}

		mockRepo.On("Search", ctx, &userID, query, limit, offset).Return(expected, nil).Once()

		result, err := service.SearchExercises(ctx, &userID, query, limit, offset)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetExercise(t *testing.T) {
	mockRepo := new(repository.MockExerciseRepository)
	service := NewExerciseService(mockRepo)
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		id := "ex-1"
		expected := &model.UniqueExercise{ID: id, Name: "Squat"}

		mockRepo.On("FindByID", ctx, id).Return(expected, nil).Once()

		result, err := service.GetExercise(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		id := "ex-999"

		mockRepo.On("FindByID", ctx, id).Return(nil, errors.New("not found")).Once()

		result, err := service.GetExercise(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
