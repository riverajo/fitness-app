package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/crypto/bcrypt" // 💡 Hashing library

	"github.com/riverajo/fitness-app/backend/internal/db"
	"github.com/riverajo/fitness-app/backend/internal/model"
)

// UserService handles all user-related business logic and database interaction.
type UserService struct {
	collection *mongo.Collection
}

// NewUserService creates a new instance of the UserService.
func NewUserService() *UserService {
	// Get the users collection on initialization
	return &UserService{
		collection: db.GetCollection("users"),
	}
}

// -------------------------------------------------------------------
// CORE BUSINESS LOGIC
// -------------------------------------------------------------------

// It assumes the user object passed in is fully prepared (hashed, timed, etc.).
func (s *UserService) CreateUser(ctx context.Context, input model.User) error {

	// Check if user already exists
	// We create a check model just for the query
	filter := primitive.M{"email": input.Email}
	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("database error during user check: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("user with email %s already exists", input.Email)
	}

	_, err = s.collection.InsertOne(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to insert user into database: %w", err)
	}

	return nil
}

// Placeholder for verification logic (needed for login)
func (s *UserService) VerifyPassword(ctx context.Context, email, password string) (*model.User, error) {
	var user model.User

	// 1. Find user by email
	filter := primitive.M{"email": email}
	err := s.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("invalid credentials")
	} else if err != nil {
		return nil, fmt.Errorf("database error during login: %w", err)
	}

	// 2. Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials") // Use a generic error message for security
	}

	return &user, nil
}

// GetUserByID fetches a user from the database using their ObjectID.
func (s *UserService) GetUserByID(ctx context.Context, idString string) (*model.User, error) {
	var user model.User

	// 1. Convert the string ID (from JWT) into a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// 2. Query MongoDB by the "_id" field
	filter := primitive.M{"_id": objectID}

	// Find one document and decode the result into the user struct
	err = s.collection.FindOne(ctx, filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, fmt.Errorf("database error fetching user: %w", err)
	}

	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, input model.UserUpdateInput) (*model.User, error) {
	// 1. Fetch the existing user (to verify current password)
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user authentication failed: %w", err)
	}

	// 2. Verify the current password (MUST BE PROVIDED)
	if input.CurrentPassword == nil || *input.CurrentPassword == "" {
		return nil, fmt.Errorf("current password is required for update")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(*input.CurrentPassword))
	if err != nil {
		return nil, fmt.Errorf("invalid current password")
	}

	// 3. Prepare the update document (MongoDB $set operation)
	update := primitive.M{}

	// Handle New Password Update
	if input.NewPassword != nil && *input.NewPassword != "" {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(*input.NewPassword), bcrypt.DefaultCost)
		if hashErr != nil {
			return nil, fmt.Errorf("failed to hash new password: %w", hashErr)
		}
		update["passwordHash"] = string(hashedPassword)
		user.PasswordHash = string(hashedPassword) // Update local model for return
	}

	// Handle Preferred Unit Update
	if input.PreferredUnit != nil && *input.PreferredUnit != "" {
		update["preferredUnit"] = *input.PreferredUnit
		user.PreferredUnit = *input.PreferredUnit // Update local model for return
	}

	// If there are updates to apply
	if len(update) > 0 {
		// Add updated timestamp
		update["updatedAt"] = time.Now()
		user.UpdatedAt = time.Now()

		_, err = s.collection.UpdateByID(ctx, user.ID, primitive.M{"$set": update})
		if err != nil {
			return nil, fmt.Errorf("database error during update: %w", err)
		}
	}

	// 4. Return the updated user entity (the local copy)
	return user, nil
}
