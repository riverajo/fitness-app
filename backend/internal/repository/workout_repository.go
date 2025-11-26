package repository

import (
	"context"

	"github.com/riverajo/fitness-app/backend/internal/model"
)

// WorkoutRepository defines the interface for workout data access.
type WorkoutRepository interface {
	Create(ctx context.Context, log model.WorkoutLog) (*model.WorkoutLog, error)
	GetByID(ctx context.Context, id string) (*model.WorkoutLog, error)
	ListByUser(ctx context.Context, userID string, limit, offset int) ([]*model.WorkoutLog, error)
	Update(ctx context.Context, log model.WorkoutLog) (*model.WorkoutLog, error)
}
