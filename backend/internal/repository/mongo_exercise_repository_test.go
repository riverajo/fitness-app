package repository

import (
	"context"
	"testing"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoExerciseRepository_Create(t *testing.T) {
	cleanupCollection(t, "unique_exercises")
	cleanupCollection(t, "unique_exercises")
	repo := NewMongoExerciseRepository(testDB)
	ctx := context.Background()

	desc := "A test exercise"
	exercise := &model.UniqueExercise{
		ID:          primitive.NewObjectID().Hex(),
		Name:        "Bench Press",
		Description: &desc,
	}

	err := repo.Create(ctx, exercise)
	assert.NoError(t, err)

	foundExercise, err := repo.FindByID(ctx, exercise.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundExercise)
	assert.Equal(t, exercise.Name, foundExercise.Name)
	assert.Equal(t, *exercise.Description, *foundExercise.Description)
}

func TestMongoExerciseRepository_FindByID(t *testing.T) {
	cleanupCollection(t, "unique_exercises")
	cleanupCollection(t, "unique_exercises")
	repo := NewMongoExerciseRepository(testDB)
	ctx := context.Background()

	exercise := &model.UniqueExercise{
		ID:   primitive.NewObjectID().Hex(),
		Name: "Squat",
	}

	err := repo.Create(ctx, exercise)
	assert.NoError(t, err)

	foundExercise, err := repo.FindByID(ctx, exercise.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundExercise)
	assert.Equal(t, exercise.Name, foundExercise.Name)

	notFoundExercise, err := repo.FindByID(ctx, primitive.NewObjectID().Hex())
	assert.NoError(t, err)
	assert.Nil(t, notFoundExercise)
}

func TestMongoExerciseRepository_Search(t *testing.T) {
	cleanupCollection(t, "unique_exercises")
	cleanupCollection(t, "unique_exercises")
	repo := NewMongoExerciseRepository(testDB)
	ctx := context.Background()

	userID := "user123"
	otherUserID := "user456"

	// System exercise (no userID)
	sysEx := &model.UniqueExercise{
		ID:   primitive.NewObjectID().Hex(),
		Name: "Push Up",
	}
	// User specific exercise
	userEx := &model.UniqueExercise{
		ID:     primitive.NewObjectID().Hex(),
		Name:   "Push Press",
		UserID: &userID,
	}
	// Other user exercise
	otherEx := &model.UniqueExercise{
		ID:     primitive.NewObjectID().Hex(),
		Name:   "Push Jerk",
		UserID: &otherUserID,
	}

	assert.NoError(t, repo.Create(ctx, sysEx))
	assert.NoError(t, repo.Create(ctx, userEx))
	assert.NoError(t, repo.Create(ctx, otherEx))

	// Search for "Push" as user123
	// Should find "Push Up" (system) and "Push Press" (user123)
	// Should NOT find "Push Jerk" (user456)
	results, err := repo.Search(ctx, &userID, "Push")
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	names := make(map[string]bool)
	for _, r := range results {
		names[r.Name] = true
	}
	assert.True(t, names["Push Up"])
	assert.True(t, names["Push Press"])
	assert.False(t, names["Push Jerk"])

	// Search for "Push" as anonymous (nil userID)
	// Should only find "Push Up"
	results, err = repo.Search(ctx, nil, "Push")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Push Up", results[0].Name)
}
