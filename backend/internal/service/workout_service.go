package service

import (
	"context"

	"github.com/riverajo/fitness-app/backend/internal/db"
	"github.com/riverajo/fitness-app/backend/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// WorkoutService defines the methods for interacting with workout data.
type WorkoutService struct {
	collection *mongo.Collection
}

// NewWorkoutService creates a new instance of the WorkoutService.
func NewWorkoutService() *WorkoutService {
	return &WorkoutService{
		collection: db.GetCollection("workout_logs"),
	}
}

// CreateLog saves a new WorkoutLog to the database.
func (s *WorkoutService) CreateLog(ctx context.Context, log model.WorkoutLog) (*model.WorkoutLog, error) {

	// 💡 1. Generate a new unique ID for the log
	log.ID = primitive.NewObjectID().Hex()

	// 💡 2. Insert the document into the collection
	result, err := s.collection.InsertOne(ctx, log)
	if err != nil {
		return nil, err
	}

	log.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return &log, nil
}
