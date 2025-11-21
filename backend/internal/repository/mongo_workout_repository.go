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

func (r *MongoWorkoutRepository) GetByID(ctx context.Context, id string) (*model.WorkoutLog, error) {
	var log model.WorkoutLog
	// Handle both string ID and ObjectID if needed, but model uses string ID mapped to _id
	// If _id is ObjectID in DB, we might need to convert.
	// Assuming _id is stored as ObjectID but we query with hex string.
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	err = r.collection.FindOne(ctx, primitive.M{"_id": oid}).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("workout log not found")
		}
		return nil, fmt.Errorf("failed to fetch workout log: %w", err)
	}
	return &log, nil
}

func (r *MongoWorkoutRepository) ListByUser(ctx context.Context, userID string) ([]*model.WorkoutLog, error) {
	cursor, err := r.collection.Find(ctx, primitive.M{"userId": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to list workout logs: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*model.WorkoutLog
	for cursor.Next(ctx) {
		var log model.WorkoutLog
		if err := cursor.Decode(&log); err != nil {
			return nil, fmt.Errorf("failed to decode workout log: %w", err)
		}
		logs = append(logs, &log)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return logs, nil
}
