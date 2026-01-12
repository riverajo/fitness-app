package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	// Use string for ID to match GraphQL and Auto-Bind.
	// We will handle ObjectID conversion in the repository.
	ID string `json:"id" bson:"_id,omitempty"`

	Email        string    `json:"email" bson:"email"`
	PasswordHash string    `json:"-" bson:"passwordHash"` // Hide from JSON/GraphQL
	CreatedAt    time.Time `json:"-" bson:"createdAt"`
	UpdatedAt    time.Time `json:"-" bson:"updatedAt"`

	// Add other internal fields
	PreferredUnit WeightUnit `json:"preferredUnit" bson:"preferredUnit"` // e.g., "KILOGRAMS" or "POUNDS"
}

// UserUpdateInput represents the fields provided for a user update.
type UserUpdateInput struct {
	CurrentPassword *string
	NewPassword     *string
	PreferredUnit   *WeightUnit
	// ... add any other updatable fields here
}

// NewUser creates a new internal User model.
func NewUser(email, hashedPassword string) *User {
	now := time.Now()

	return &User{
		ID:            primitive.NewObjectID().Hex(), // Generate DB ID
		Email:         email,
		PasswordHash:  hashedPassword, // Takes the already-hashed password
		CreatedAt:     now,
		UpdatedAt:     now,
		PreferredUnit: WeightUnitKilograms, // Set default value
	}
}
