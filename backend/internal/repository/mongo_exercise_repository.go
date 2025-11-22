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

type MongoExerciseRepository struct {
	collection *mongo.Collection
}

func NewMongoExerciseRepository(database *mongo.Database) *MongoExerciseRepository {
	return &MongoExerciseRepository{
		collection: database.Collection("unique_exercises"),
	}
}

func (r *MongoExerciseRepository) Create(ctx context.Context, exercise *model.UniqueExercise) error {
	// Ensure ID is generated if empty
	if exercise.ID == "" {
		exercise.ID = primitive.NewObjectID().Hex()
	}

	oid, err := primitive.ObjectIDFromHex(exercise.ID)
	if err != nil {
		return fmt.Errorf("invalid exercise ID format: %w", err)
	}

	doc := bson.M{
		"_id":  oid,
		"name": exercise.Name,
	}
	if exercise.UserID != nil {
		doc["userId"] = *exercise.UserID
	}
	if exercise.Description != nil {
		doc["description"] = *exercise.Description
	}

	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to insert exercise: %w", err)
	}

	return nil
}

func (r *MongoExerciseRepository) Search(ctx context.Context, userID *string, query string) ([]*model.UniqueExercise, error) {
	// Filter: (userId == nil OR userId == currentUserId) AND name matches query
	filter := bson.M{
		"name": bson.M{"$regex": query, "$options": "i"},
	}

	userFilter := bson.A{bson.M{"userId": nil}} // System exercises
	if userID != nil {
		userFilter = append(userFilter, bson.M{"userId": *userID})
	}
	filter["$or"] = userFilter

	// Limit results to prevent overload
	opts := options.Find().SetLimit(50)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("database error searching exercises: %w", err)
	}
	defer cursor.Close(ctx)

	var exercises []*model.UniqueExercise
	for cursor.Next(ctx) {
		var doc struct {
			ID          primitive.ObjectID `bson:"_id"`
			Name        string             `bson:"name"`
			UserID      *string            `bson:"userId"`
			Description *string            `bson:"description"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode exercise: %w", err)
		}

		exercises = append(exercises, &model.UniqueExercise{
			ID:          doc.ID.Hex(),
			Name:        doc.Name,
			UserID:      doc.UserID,
			Description: doc.Description,
		})
	}

	return exercises, nil
}

func (r *MongoExerciseRepository) FindByID(ctx context.Context, id string) (*model.UniqueExercise, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid exercise ID format: %w", err)
	}

	var doc struct {
		ID          primitive.ObjectID `bson:"_id"`
		Name        string             `bson:"name"`
		UserID      *string            `bson:"userId"`
		Description *string            `bson:"description"`
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("database error finding exercise by ID: %w", err)
	}

	return &model.UniqueExercise{
		ID:          doc.ID.Hex(),
		Name:        doc.Name,
		UserID:      doc.UserID,
		Description: doc.Description,
	}, nil
}
