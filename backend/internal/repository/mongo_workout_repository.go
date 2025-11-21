package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
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
	if log.ID == "" {
		log.ID = primitive.NewObjectID().Hex()
	}

	oid, err := primitive.ObjectIDFromHex(log.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	// Create a separate struct or map to ensure _id is inserted as ObjectID
	// We can use a map to avoid defining a new struct if we trust the model fields match bson tags
	// But model.WorkoutLog has `bson:"_id,omitempty"` on ID string.
	// So we can't just pass log.
	// We'll use a map for simplicity or a struct wrapper.
	// Given the model is complex (nested slices), map might be easier but we need to be careful with types.
	// Actually, we can just use the model but overwrite _id if we could, but we can't change type of ID field.
	// So we must marshal to something else.

	// Let's use a map, but we need to copy all fields. That's tedious and error prone.
	// Better approach: Define an alias or struct with ObjectID _id.
	// Or, since we are in the repository, we can define a private struct that mirrors the model but with ObjectID.
	// But the model has nested structs (ExerciseLog, Set) which also have bson tags.
	// If we use a map, we can marshal the log to bytes then unmarshal to map, then fix _id?
	// That's slow.

	// Let's try to just use a map for the top level fields.
	doc := bson.M{
		"_id":          oid,
		"userId":       log.UserID,
		"name":         log.Name,
		"startTime":    log.StartTime,
		"endTime":      log.EndTime,
		"exerciseLogs": log.ExerciseLogs,
		"locationName": log.LocationName,
		"generalNotes": log.GeneralNotes,
	}

	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to insert workout log: %w", err)
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
