package repository

import (
	"context"
	"testing"
	"time"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoWorkoutRepository_Create(t *testing.T) {
	cleanupCollection(t, "workout_logs")
	cleanupCollection(t, "workout_logs")
	repo := NewMongoWorkoutRepository(testDB)
	ctx := context.Background()

	workout := model.WorkoutLog{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    primitive.NewObjectID().Hex(),
		StartTime: time.Now(),
		EndTime:   time.Now().Add(1 * time.Hour),
		Name:      "Morning Workout",
	}

	createdWorkout, err := repo.Create(ctx, workout)
	require.NoError(t, err)
	require.NotNil(t, createdWorkout)
	assert.Equal(t, workout.ID, createdWorkout.ID)

	foundWorkout, err := repo.GetByID(ctx, workout.ID)
	require.NoError(t, err)
	require.NotNil(t, foundWorkout)
	assert.Equal(t, workout.ID, foundWorkout.ID)
	assert.Equal(t, workout.UserID, foundWorkout.UserID)
}

func TestMongoWorkoutRepository_GetByID(t *testing.T) {
	cleanupCollection(t, "workout_logs")
	cleanupCollection(t, "workout_logs")
	repo := NewMongoWorkoutRepository(testDB)
	ctx := context.Background()

	workout := model.WorkoutLog{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    primitive.NewObjectID().Hex(),
		StartTime: time.Now(),
		EndTime:   time.Now().Add(45 * time.Minute),
	}

	_, err := repo.Create(ctx, workout)
	require.NoError(t, err)

	foundWorkout, err := repo.GetByID(ctx, workout.ID)
	require.NoError(t, err)
	require.NotNil(t, foundWorkout)
	assert.Equal(t, workout.ID, foundWorkout.ID)

	// Note: GetByID returns error if not found, not nil
	_, err = repo.GetByID(ctx, primitive.NewObjectID().Hex())
	assert.Error(t, err)
}

func TestMongoWorkoutRepository_ListByUser(t *testing.T) {
	cleanupCollection(t, "workout_logs")
	cleanupCollection(t, "workout_logs")
	repo := NewMongoWorkoutRepository(testDB)
	ctx := context.Background()

	userID := primitive.NewObjectID().Hex()

	workout1 := model.WorkoutLog{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    userID,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(30 * time.Minute),
	}
	workout2 := model.WorkoutLog{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    userID,
		StartTime: time.Now().Add(24 * time.Hour),
		EndTime:   time.Now().Add(25 * time.Hour),
	}
	workout3 := model.WorkoutLog{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    primitive.NewObjectID().Hex(), // Different user
		StartTime: time.Now(),
		EndTime:   time.Now().Add(45 * time.Minute),
	}

	_, err := repo.Create(ctx, workout1)
	require.NoError(t, err)
	_, err = repo.Create(ctx, workout2)
	require.NoError(t, err)
	_, err = repo.Create(ctx, workout3)
	require.NoError(t, err)

	workouts, err := repo.ListByUser(ctx, userID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, workouts, 2)

	// Verify we got the correct workouts
	ids := make(map[string]bool)
	for _, w := range workouts {
		ids[w.ID] = true
	}
	assert.True(t, ids[workout1.ID])
	assert.True(t, ids[workout2.ID])
	assert.False(t, ids[workout3.ID])
}
