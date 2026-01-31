package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	FindByID(ctx context.Context, id string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllForUser(ctx context.Context, userID string) error
}

type MongoRefreshTokenRepository struct {
	collection *mongo.Collection
}

func NewMongoRefreshTokenRepository(db *mongo.Database) *MongoRefreshTokenRepository {
	collection := db.Collection("refreshtokens")

	// Create TTL Index on ExpiresAt
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	// Use a background context for index creation as it happens on startup
	if _, err := collection.Indexes().CreateOne(context.Background(), indexModel); err != nil {
		slog.Error("Failed to create TTL index for refresh tokens", "error", err)
	}

	return &MongoRefreshTokenRepository{
		collection: collection,
	}
}

func (r *MongoRefreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	if token.ID == "" {
		token.ID = bson.NewObjectID().Hex()
	}
	_, err := r.collection.InsertOne(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

func (r *MongoRefreshTokenRepository) FindByID(ctx context.Context, id string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	filter := bson.D{{Key: "_id", Value: id}}
	err := r.collection.FindOne(ctx, filter).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}
	return &token, nil
}

func (r *MongoRefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

func (r *MongoRefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	filter := bson.D{{Key: "userId", Value: userID}}
	_, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to revoke all tokens for user: %w", err)
	}
	return nil
}
