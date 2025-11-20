package service

import (
	"context"
	"errors"
	"testing"

	"github.com/riverajo/fitness-app/backend/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	CreateFunc      func(ctx context.Context, user model.User) error
	FindByEmailFunc func(ctx context.Context, email string) (*model.User, error)
	FindByIDFunc    func(ctx context.Context, id string) (*model.User, error)
	UpdateFunc      func(ctx context.Context, user *model.User) error
}

func (m *MockUserRepository) Create(ctx context.Context, user model.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name    string
		input   model.User
		mock    *MockUserRepository
		wantErr bool
	}{
		{
			name: "Success",
			input: model.User{
				Email: "test@example.com",
			},
			mock: &MockUserRepository{
				CreateFunc: func(ctx context.Context, user model.User) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "Error from Repository",
			input: model.User{
				Email: "test@example.com",
			},
			mock: &MockUserRepository{
				CreateFunc: func(ctx context.Context, user model.User) error {
					return errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(tt.mock)
			err := s.CreateUser(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name     string
		email    string
		password string
		mock     *MockUserRepository
		wantUser bool
		wantErr  bool
	}{
		{
			name:     "Success",
			email:    "test@example.com",
			password: "password123",
			mock: &MockUserRepository{
				FindByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
					return &model.User{
						Email:        "test@example.com",
						PasswordHash: string(hashedPassword),
					}, nil
				},
			},
			wantUser: true,
			wantErr:  false,
		},
		{
			name:     "User Not Found",
			email:    "unknown@example.com",
			password: "password123",
			mock: &MockUserRepository{
				FindByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
					return nil, nil
				},
			},
			wantUser: false,
			wantErr:  true,
		},
		{
			name:     "Wrong Password",
			email:    "test@example.com",
			password: "wrongpassword",
			mock: &MockUserRepository{
				FindByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
					return &model.User{
						Email:        "test@example.com",
						PasswordHash: string(hashedPassword),
					}, nil
				},
			},
			wantUser: false,
			wantErr:  true,
		},
		{
			name:     "DB Error",
			email:    "test@example.com",
			password: "password123",
			mock: &MockUserRepository{
				FindByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
					return nil, errors.New("db error")
				},
			},
			wantUser: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(tt.mock)
			user, err := s.VerifyPassword(context.Background(), tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantUser && user == nil {
				t.Error("VerifyPassword() expected user, got nil")
			}
			if !tt.wantUser && user != nil {
				t.Error("VerifyPassword() expected nil user, got user")
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	validID := primitive.NewObjectID()
	tests := []struct {
		name     string
		id       string
		mock     *MockUserRepository
		wantUser bool
		wantErr  bool
	}{
		{
			name: "Success",
			id:   validID.Hex(),
			mock: &MockUserRepository{
				FindByIDFunc: func(ctx context.Context, id string) (*model.User, error) {
					return &model.User{ID: validID}, nil
				},
			},
			wantUser: true,
			wantErr:  false,
		},
		{
			name: "Not Found",
			id:   validID.Hex(),
			mock: &MockUserRepository{
				FindByIDFunc: func(ctx context.Context, id string) (*model.User, error) {
					return nil, nil
				},
			},
			wantUser: false,
			wantErr:  true,
		},
		{
			name: "DB Error",
			id:   validID.Hex(),
			mock: &MockUserRepository{
				FindByIDFunc: func(ctx context.Context, id string) (*model.User, error) {
					return nil, errors.New("db error")
				},
			},
			wantUser: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserService(tt.mock)
			user, err := s.GetUserByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantUser && user == nil {
				t.Error("GetUserByID() expected user, got nil")
			}
		})
	}
}
