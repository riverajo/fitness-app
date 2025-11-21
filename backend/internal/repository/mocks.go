package repository

import (
	"context"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockWorkoutRepository is a mock implementation of WorkoutRepository
type MockWorkoutRepository struct {
	mock.Mock
}

func (m *MockWorkoutRepository) Create(ctx context.Context, log model.WorkoutLog) (*model.WorkoutLog, error) {
	args := m.Called(ctx, log)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WorkoutLog), args.Error(1)
}

func (m *MockWorkoutRepository) GetByID(ctx context.Context, id string) (*model.WorkoutLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WorkoutLog), args.Error(1)
}

func (m *MockWorkoutRepository) ListByUser(ctx context.Context, userID string) ([]*model.WorkoutLog, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.WorkoutLog), args.Error(1)
}
