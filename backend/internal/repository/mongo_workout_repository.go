package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/riverajo/fitness-app/backend/internal/db"
	"github.com/riverajo/fitness-app/backend/internal/model"
)

type MongoWorkoutRepository struct {
	collection *mongo.Collection
}

func NewMongoWorkoutRepository() *MongoWorkoutRepository {
	return &MongoWorkoutRepository{
		collection: db.GetCollection("workout_logs"),
	}
}

func (r *MongoWorkoutRepository) Create(ctx context.Context, log model.WorkoutLog) (*model.WorkoutLog, error) {
	// Generate a new unique ID for the log if not present (though model might have string ID)
	// The service was doing: log.ID = primitive.NewObjectID().Hex()
	// We should probably stick to that or let Mongo generate it.
	// If we use string IDs in the model, we need to handle that.

	if log.ID == "" {
		log.ID = primitive.NewObjectID().Hex()
	}

	result, err := r.collection.InsertOne(ctx, log)
	if err != nil {
		return nil, fmt.Errorf("failed to insert workout log: %w", err)
	}

	// Ensure the ID is set correctly on the returned object
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		log.ID = oid.Hex()
	} else if idStr, ok := result.InsertedID.(string); ok {
		log.ID = idStr
	}

	return &log, nil
}
