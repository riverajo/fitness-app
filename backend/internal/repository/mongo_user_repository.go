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

	// Convert string ID to ObjectID for storage if it's a valid hex string
	// However, if we store it as _id, we need to be careful.
	// The model has `ID string `bson:"_id,omitempty"``.
	// If we pass the struct directly to InsertOne, the driver will try to insert the string as _id.
	// But we want _id to be an ObjectID.
	// So we need to map it to a struct that uses ObjectID or let the driver generate it if empty.
	// But NewUserFromRegisterInput generates a hex string.
	// We should parse it back to ObjectID for insertion.

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
