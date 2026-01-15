package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/riverajo/fitness-app/backend/internal/model"
)

type MongoWorkoutRepository struct {
	collection *mongo.Collection
}

func NewMongoWorkoutRepository(database *mongo.Database) *MongoWorkoutRepository {
	return &MongoWorkoutRepository{
		collection: database.Collection("workout_logs"),
	}
}

func (r *MongoWorkoutRepository) Create(ctx context.Context, logData model.WorkoutLog) (*model.WorkoutLog, error) {
	if logData.ID == "" {
		logData.ID = primitive.NewObjectID().Hex()
	}

	oid, err := primitive.ObjectIDFromHex(logData.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	// Use a map to ensure _id is inserted as an ObjectID rather than a string.
	// We map the model fields manually to the BSON map.
	doc := bson.M{
		"_id":          oid,
		"userId":       logData.UserID,
		"name":         logData.Name,
		"startTime":    logData.StartTime,
		"endTime":      logData.EndTime,
		"exerciseLogs": logData.ExerciseLogs,
		"locationName": logData.LocationName,
		"generalNotes": logData.GeneralNotes,
	}

	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to insert workout log: %w", err)
	}

	return &logData, nil
}

func (r *MongoWorkoutRepository) GetByID(ctx context.Context, id string) (*model.WorkoutLog, error) {
	var log model.WorkoutLog
	// Convert hex string ID to ObjectID.
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

func (r *MongoWorkoutRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*model.WorkoutLog, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetSort(bson.D{{Key: "startTime", Value: -1}})
	cursor, err := r.collection.Find(ctx, primitive.M{"userId": userID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list workout logs: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

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

func (r *MongoWorkoutRepository) Update(ctx context.Context, logData model.WorkoutLog) (*model.WorkoutLog, error) {
	oid, err := primitive.ObjectIDFromHex(logData.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"name":         logData.Name,
			"startTime":    logData.StartTime,
			"endTime":      logData.EndTime,
			"exerciseLogs": logData.ExerciseLogs,
			"locationName": logData.LocationName,
			"generalNotes": logData.GeneralNotes,
		},
	}

	// Filter by _id and optionally userId to ensure ownership.
	filter := primitive.M{"_id": oid}
	if logData.UserID != "" {
		filter["userId"] = logData.UserID
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update workout log: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("workout log not found or unauthorized")
	}

	return &logData, nil
}
