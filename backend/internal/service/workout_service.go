package service

import (
	"context"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/riverajo/fitness-app/backend/internal/repository"
)

// WorkoutService defines the methods for interacting with workout data.
type WorkoutService struct {
	repo repository.WorkoutRepository
}

// NewWorkoutService creates a new instance of the WorkoutService.
func NewWorkoutService(repo repository.WorkoutRepository) *WorkoutService {
	return &WorkoutService{
		repo: repo,
	}
}

// CreateLog saves a new WorkoutLog to the database.
func (s *WorkoutService) CreateLog(ctx context.Context, log model.WorkoutLog) (*model.WorkoutLog, error) {
	return s.repo.Create(ctx, log)
}
