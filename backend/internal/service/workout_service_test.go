package service

import (
	"context"
	"errors"
	"testing"

	"github.com/riverajo/fitness-app/backend/internal/model"
)

// MockWorkoutRepository is a mock implementation of repository.WorkoutRepository
type MockWorkoutRepository struct {
	CreateFunc     func(ctx context.Context, log model.WorkoutLog) (*model.WorkoutLog, error)
	GetByIDFunc    func(ctx context.Context, id string) (*model.WorkoutLog, error)
	ListByUserFunc func(ctx context.Context, userID string, limit, offset int) ([]*model.WorkoutLog, error)
}

func (m *MockWorkoutRepository) Create(ctx context.Context, log model.WorkoutLog) (*model.WorkoutLog, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, log)
	}
	return &log, nil
}

func (m *MockWorkoutRepository) GetByID(ctx context.Context, id string) (*model.WorkoutLog, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockWorkoutRepository) ListByUser(ctx context.Context, userID string, limit, offset int) ([]*model.WorkoutLog, error) {
	if m.ListByUserFunc != nil {
		return m.ListByUserFunc(ctx, userID, limit, offset)
	}
	return nil, nil
}

func TestCreateLog(t *testing.T) {
	tests := []struct {
		name    string
		input   model.WorkoutLog
		mock    *MockWorkoutRepository
		wantErr bool
	}{
		{
			name: "Success",
			input: model.WorkoutLog{
				Name: "Leg Day",
			},
			mock: &MockWorkoutRepository{
				CreateFunc: func(ctx context.Context, log model.WorkoutLog) (*model.WorkoutLog, error) {
					log.ID = "generated-id"
					return &log, nil
				},
			},
			wantErr: false,
		},
		{
			name: "Repository Error",
			input: model.WorkoutLog{
				Name: "Leg Day",
			},
			mock: &MockWorkoutRepository{
				CreateFunc: func(ctx context.Context, log model.WorkoutLog) (*model.WorkoutLog, error) {
					return nil, errors.New("db error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewWorkoutService(tt.mock)
			got, err := s.CreateLog(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("CreateLog() expected result, got nil")
				} else if got.ID == "" {
					t.Error("CreateLog() expected ID to be generated")
				}
			}
		})
	}
}
