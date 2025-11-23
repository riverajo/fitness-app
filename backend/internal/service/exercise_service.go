package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/riverajo/fitness-app/backend/internal/repository"
)

type ExerciseService struct {
	repo repository.ExerciseRepository
}

func NewExerciseService(repo repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{
		repo: repo,
	}
}

func (s *ExerciseService) CreateExercise(ctx context.Context, name string, description *string, userID *string) (*model.UniqueExercise, error) {
	// 1. Validate input
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("exercise name cannot be empty")
	}

	// 2. Check for duplicates (optional but good practice)
	// For now, we'll skip a strict duplicate check to keep it simple,
	// or we could search for exact matches.
	// Let's trust the user for now, or maybe just check if they already have one with that exact name?
	// The search logic is fuzzy, so it might be hard to check exact equality efficiently without a specific query.
	// Let's proceed with creation.

	exercise := &model.UniqueExercise{
		Name:        name,
		UserID:      userID,
		Description: description,
	}

	if err := s.repo.Create(ctx, exercise); err != nil {
		return nil, err
	}

	return exercise, nil
}

func (s *ExerciseService) SearchExercises(ctx context.Context, userID *string, query string, limit int, offset int) ([]*model.UniqueExercise, error) {
	return s.repo.Search(ctx, userID, query, limit, offset)
}

func (s *ExerciseService) GetExercise(ctx context.Context, id string) (*model.UniqueExercise, error) {
	return s.repo.FindByID(ctx, id)
}
