package repository

import (
	"context"

	"github.com/riverajo/fitness-app/backend/internal/model"
)

// WorkoutRepository defines the interface for workout data access.
type WorkoutRepository interface {
	Create(ctx context.Context, log model.WorkoutLog) (*model.WorkoutLog, error)
}
