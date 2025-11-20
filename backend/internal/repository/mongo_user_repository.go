package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/riverajo/fitness-app/backend/internal/db"
	"github.com/riverajo/fitness-app/backend/internal/model"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository() *MongoUserRepository {
	return &MongoUserRepository{
		collection: db.GetCollection("users"),
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

	_, err = r.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert user into database: %w", err)
	}

	return nil
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	filter := primitive.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil // Return nil if not found, let service handle error
	} else if err != nil {
		return nil, fmt.Errorf("database error finding user by email: %w", err)
	}
	return &user, nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	filter := primitive.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil // Return nil if not found
	} else if err != nil {
		return nil, fmt.Errorf("database error finding user by ID: %w", err)
	}
	return &user, nil
}

func (r *MongoUserRepository) Update(ctx context.Context, user *model.User) error {
	// We only update specific fields that are allowed to change
	// Ideally, this should be more granular, but for now we replicate existing logic
	// actually the existing logic constructs the update map dynamically.
	// Let's stick to the interface which takes a *model.User and updates it.
	// However, the service logic was doing partial updates.
	// The interface I defined `Update(ctx context.Context, user *model.User) error` implies saving the state of the user.
	// But the service logic was selectively updating fields.

	// Let's refine the implementation to match the service's needs or adjust the service.
	// The service constructs a `$set` map.
	// If we want to keep the repository generic, `Update` usually saves the whole entity or specific fields.
	// Let's implement a full update for the fields that can change, or we can just update what's in the struct.

	// Re-reading the service logic:
	// It checks input.NewPassword and input.PreferredUnit and adds them to `update` map.
	// Then it calls UpdateByID.

	// To support this in the repository without leaking "UserUpdateInput" (which is a GQL model) into the repository (which deals with domain models),
	// passing the updated `model.User` struct is correct. The service should prepare the `model.User` with the new values.
	// But `model.User` has all fields. If we just save `model.User`, we might overwrite things we didn't intend to if the struct isn't fully populated?
	// No, `FindByID` returns a full struct. The service modifies it. So passing it back to `Update` is safe if we update all fields or just the ones we care about.

	// Let's update the fields that are mutable.

	updateFields := primitive.M{
		"passwordHash":  user.PasswordHash,
		"preferredUnit": user.PreferredUnit,
		"updatedAt":     time.Now(),
	}

	_, err := r.collection.UpdateByID(ctx, user.ID, primitive.M{"$set": updateFields})
	if err != nil {
		return fmt.Errorf("database error during update: %w", err)
	}
	return nil
}
