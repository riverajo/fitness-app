package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/riverajo/fitness-app/backend/internal/service" // 💡 New import
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	// 💡 Inject the service dependencies here
	WorkoutService *service.WorkoutService
	UserService    *service.UserService // For authentication/user CRUD
	// Add other services (e.g., ExerciseService) as you create them
}

func NewResolver() *Resolver {
	return &Resolver{
		UserService:    service.NewUserService(),
		WorkoutService: service.NewWorkoutService(),
	}
}
