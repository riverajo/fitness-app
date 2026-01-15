package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/riverajo/fitness-app/backend/internal/model"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(database *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		collection: database.Collection("users"),
	}
}

func (r *MongoUserRepository) Create(ctx context.Context, user model.User) error {
	// Check if user already exists
	filter := primitive.M{"email": user.Email}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("database error during user check: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Convert hex string ID to ObjectID for storage.

	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Create a temporary struct or map for insertion to ensure _id is ObjectID
	userDoc := struct {
		ID            primitive.ObjectID `bson:"_id,omitempty"`
		Email         string             `bson:"email"`
		PasswordHash  string             `bson:"passwordHash"`
		CreatedAt     time.Time          `bson:"createdAt"`
		UpdatedAt     time.Time          `bson:"updatedAt"`
		PreferredUnit string             `bson:"preferredUnit"`
	}{
		ID:            oid,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		PreferredUnit: string(user.PreferredUnit),
	}

	_, err = r.collection.InsertOne(ctx, userDoc)
	if err != nil {
		return fmt.Errorf("failed to insert user into database: %w", err)
	}

	return nil
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	// We need to decode into a struct that matches the DB (ObjectID _id) and then map to model.User (string ID)
	var userDoc struct {
		ID            primitive.ObjectID `bson:"_id"`
		Email         string             `bson:"email"`
		PasswordHash  string             `bson:"passwordHash"`
		CreatedAt     time.Time          `bson:"createdAt"`
		UpdatedAt     time.Time          `bson:"updatedAt"`
		PreferredUnit string             `bson:"preferredUnit"`
	}

	filter := primitive.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&userDoc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("database error finding user by email: %w", err)
	}

	return &model.User{
		ID:            userDoc.ID.Hex(),
		Email:         userDoc.Email,
		PasswordHash:  userDoc.PasswordHash,
		CreatedAt:     userDoc.CreatedAt,
		UpdatedAt:     userDoc.UpdatedAt,
		PreferredUnit: model.WeightUnit(userDoc.PreferredUnit),
	}, nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	var userDoc struct {
		ID            primitive.ObjectID `bson:"_id"`
		Email         string             `bson:"email"`
		PasswordHash  string             `bson:"passwordHash"`
		CreatedAt     time.Time          `bson:"createdAt"`
		UpdatedAt     time.Time          `bson:"updatedAt"`
		PreferredUnit string             `bson:"preferredUnit"`
	}

	filter := primitive.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&userDoc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("database error finding user by ID: %w", err)
	}

	return &model.User{
		ID:            userDoc.ID.Hex(),
		Email:         userDoc.Email,
		PasswordHash:  userDoc.PasswordHash,
		CreatedAt:     userDoc.CreatedAt,
		UpdatedAt:     userDoc.UpdatedAt,
		PreferredUnit: model.WeightUnit(userDoc.PreferredUnit),
	}, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, user *model.User) error {
	// Update mutable fields.

	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	updateFields := primitive.M{
		"passwordHash":  user.PasswordHash,
		"preferredUnit": user.PreferredUnit,
		"updatedAt":     time.Now(),
	}

	_, err = r.collection.UpdateByID(ctx, oid, primitive.M{"$set": updateFields})
	if err != nil {
		return fmt.Errorf("database error during update: %w", err)
	}
	return nil
}
