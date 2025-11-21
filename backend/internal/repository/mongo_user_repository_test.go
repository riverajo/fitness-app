package repository

import (
	"context"
	"testing"
	"time"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoUserRepository_Create(t *testing.T) {
	cleanupCollection(t, "users")
	repo := NewMongoUserRepository()
	ctx := context.Background()

	user := model.User{
		ID:            primitive.NewObjectID().Hex(),
		Email:         "test@example.com",
		PasswordHash:  "hashedpassword",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		PreferredUnit: "metric",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Verify insertion
	foundUser, err := repo.FindByEmail(ctx, user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.ID, foundUser.ID)
}

func TestMongoUserRepository_FindByEmail(t *testing.T) {
	cleanupCollection(t, "users")
	repo := NewMongoUserRepository()
	ctx := context.Background()

	user := model.User{
		ID:            primitive.NewObjectID().Hex(),
		Email:         "findme@example.com",
		PasswordHash:  "hashedpassword",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		PreferredUnit: "metric",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	foundUser, err := repo.FindByEmail(ctx, "findme@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)

	notFoundUser, err := repo.FindByEmail(ctx, "nonexistent@example.com")
	assert.NoError(t, err)
	assert.Nil(t, notFoundUser)
}

func TestMongoUserRepository_FindByID(t *testing.T) {
	cleanupCollection(t, "users")
	repo := NewMongoUserRepository()
	ctx := context.Background()

	user := model.User{
		ID:            primitive.NewObjectID().Hex(),
		Email:         "findbyid@example.com",
		PasswordHash:  "hashedpassword",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		PreferredUnit: "metric",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	foundUser, err := repo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.ID, foundUser.ID)

	notFoundUser, err := repo.FindByID(ctx, primitive.NewObjectID().Hex())
	assert.NoError(t, err)
	assert.Nil(t, notFoundUser)
}

func TestMongoUserRepository_Update(t *testing.T) {
	cleanupCollection(t, "users")
	repo := NewMongoUserRepository()
	ctx := context.Background()

	user := model.User{
		ID:            primitive.NewObjectID().Hex(),
		Email:         "update@example.com",
		PasswordHash:  "oldhash",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		PreferredUnit: "metric",
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	user.PasswordHash = "newhash"
	user.PreferredUnit = "imperial"

	err = repo.Update(ctx, &user)
	assert.NoError(t, err)

	updatedUser, err := repo.FindByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "newhash", updatedUser.PasswordHash)
	assert.Equal(t, "imperial", updatedUser.PreferredUnit)
}
