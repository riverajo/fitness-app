package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

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
		logData.ID = bson.NewObjectID().Hex()
	}

	oid, err := bson.ObjectIDFromHex(logData.ID)
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
	// Convert hex string ID to ObjectID.
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	var doc struct {
		ID           bson.ObjectID        `bson:"_id"`
		UserID       string               `bson:"userId"`
		Name         string               `bson:"name"`
		StartTime    time.Time            `bson:"startTime"`
		EndTime      time.Time            `bson:"endTime"`
		ExerciseLogs []*model.ExerciseLog `bson:"exerciseLogs"`
		LocationName *string              `bson:"locationName"`
		GeneralNotes *string              `bson:"generalNotes"`
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("workout log not found")
		}
		return nil, fmt.Errorf("failed to fetch workout log: %w", err)
	}
	return &model.WorkoutLog{
		ID:           doc.ID.Hex(),
		UserID:       doc.UserID,
		Name:         doc.Name,
		StartTime:    doc.StartTime,
		EndTime:      doc.EndTime,
		ExerciseLogs: doc.ExerciseLogs,
		LocationName: doc.LocationName,
		GeneralNotes: doc.GeneralNotes,
	}, nil
}

func (r *MongoWorkoutRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*model.WorkoutLog, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetSort(bson.D{{Key: "startTime", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list workout logs: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var logs []*model.WorkoutLog
	for cursor.Next(ctx) {
		var doc struct {
			ID           bson.ObjectID        `bson:"_id"`
			UserID       string               `bson:"userId"`
			Name         string               `bson:"name"`
			StartTime    time.Time            `bson:"startTime"`
			EndTime      time.Time            `bson:"endTime"`
			ExerciseLogs []*model.ExerciseLog `bson:"exerciseLogs"`
			LocationName *string              `bson:"locationName"`
			GeneralNotes *string              `bson:"generalNotes"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode workout log: %w", err)
		}
		logs = append(logs, &model.WorkoutLog{
			ID:           doc.ID.Hex(),
			UserID:       doc.UserID,
			Name:         doc.Name,
			StartTime:    doc.StartTime,
			EndTime:      doc.EndTime,
			ExerciseLogs: doc.ExerciseLogs,
			LocationName: doc.LocationName,
			GeneralNotes: doc.GeneralNotes,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return logs, nil
}

func (r *MongoWorkoutRepository) Update(ctx context.Context, logData model.WorkoutLog) (*model.WorkoutLog, error) {
	oid, err := bson.ObjectIDFromHex(logData.ID)
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
	filter := bson.M{"_id": oid}
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
