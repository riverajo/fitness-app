package service

import (
	"context"
	"errors"
	"testing"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/riverajo/fitness-app/backend/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateLog(t *testing.T) {
	mockRepo := new(repository.MockWorkoutRepository)
	service := NewWorkoutService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		input := model.WorkoutLog{Name: "Leg Day"}
		mockRepo.On("Create", ctx, input).Return(&model.WorkoutLog{ID: "generated-id", Name: "Leg Day"}, nil).Once()

		result, err := service.CreateLog(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, "generated-id", result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		input := model.WorkoutLog{Name: "Leg Day"}
		mockRepo.On("Create", ctx, input).Return(nil, errors.New("db error")).Once()

		result, err := service.CreateLog(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetLog(t *testing.T) {
	mockRepo := new(repository.MockWorkoutRepository)
	service := NewWorkoutService(mockRepo)
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		id := "log-1"
		expected := &model.WorkoutLog{ID: id, Name: "Upper Body"}
		mockRepo.On("GetByID", ctx, id).Return(expected, nil).Once()

		result, err := service.GetLog(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		id := "log-999"
		mockRepo.On("GetByID", ctx, id).Return(nil, errors.New("not found")).Once()

		result, err := service.GetLog(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestListLogs(t *testing.T) {
	mockRepo := new(repository.MockWorkoutRepository)
	service := NewWorkoutService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userID := "user-123"
		limit := 10
		offset := 0
		expected := []*model.WorkoutLog{{Name: "Run"}, {Name: "Gym"}}
		mockRepo.On("ListByUser", ctx, userID, limit, offset).Return(expected, nil).Once()

		result, err := service.ListLogs(ctx, userID, limit, offset)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateLog(t *testing.T) {
	mockRepo := new(repository.MockWorkoutRepository)
	service := NewWorkoutService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		input := model.WorkoutLog{ID: "log-1", Name: "Updated Name"}
		expected := &model.WorkoutLog{ID: "log-1", Name: "Updated Name"}
		mockRepo.On("Update", ctx, input).Return(expected, nil).Once()

		result, err := service.UpdateLog(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		input := model.WorkoutLog{ID: "log-1"}
		mockRepo.On("Update", ctx, input).Return(nil, errors.New("update failed")).Once()

		result, err := service.UpdateLog(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
