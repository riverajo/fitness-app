package model

import (
	"time"
)

type RefreshToken struct {
	// ID is the unique identifier for the refresh token document.
	// We use string to align with the User ID pattern, though MongoDB uses ObjectIDs naturally.
	ID string `bson:"_id,omitempty"`

	UserID    string    `bson:"userId"`
	TokenHash string    `bson:"tokenHash"`
	ExpiresAt time.Time `bson:"expiresAt"`
	CreatedAt time.Time `bson:"createdAt"`
}
