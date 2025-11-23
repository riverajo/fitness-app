package repository

import (
	"context"

	"github.com/riverajo/fitness-app/backend/internal/model"
)

type ExerciseRepository interface {
	Create(ctx context.Context, exercise *model.UniqueExercise) error
	Search(ctx context.Context, userID *string, query string, limit int, offset int) ([]*model.UniqueExercise, error)
	FindByID(ctx context.Context, id string) (*model.UniqueExercise, error)
}
