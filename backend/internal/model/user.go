package model

import (
	"time"
	// Import the generated GraphQL model package for the input type
	gqlModel "github.com/riverajo/fitness-app/backend/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	// Use primitive.ObjectID for MongoDB's internal ID, mapped to BSON "_id".
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Email        string    `bson:"email"`
	PasswordHash string    `bson:"passwordHash"` // 💡 Crucial: Stores the secure hash, NOT the plain text password.
	CreatedAt    time.Time `bson:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt"`

	// Add other internal fields
	PreferredUnit string `bson:"preferredUnit"` // e.g., "KGS" or "LBS"
}

// UserUpdateInput represents the fields provided for a user update.
type UserUpdateInput struct {
	CurrentPassword *string
	NewPassword     *string
	PreferredUnit   *string
	// ... add any other updatable fields here
}

// NewUserFromRegisterInput takes the generated GraphQL RegisterInput and
// converts it into a new internal User model, ready for service use.
func NewUserFromRegisterInput(input gqlModel.RegisterInput, hashedPassword string) *User {
	now := time.Now()

	return &User{
		ID:            primitive.NewObjectID(), // Generate DB ID
		Email:         input.Email,
		PasswordHash:  hashedPassword, // Takes the already-hashed password
		CreatedAt:     now,
		UpdatedAt:     now,
		PreferredUnit: WeightUnitKilograms, // Set default value
	}
}

// ToGraphQLUser converts the internal DB model (which contains the hash)
// into the public GraphQL model (which is safe to expose).
func (u *User) ToGraphQLUser() *gqlModel.User {

	return &gqlModel.User{
		ID:            u.ID.Hex(),
		Email:         u.Email,
		PreferredUnit: u.PreferredUnit,
	}
}
