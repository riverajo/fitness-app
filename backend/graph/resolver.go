package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/riverajo/fitness-app/backend/internal/config"
	"github.com/riverajo/fitness-app/backend/internal/repository"
	"github.com/riverajo/fitness-app/backend/internal/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	// 💡 Inject the service dependencies here
	WorkoutService  *service.WorkoutService
	UserService     *service.UserService
	ExerciseService *service.ExerciseService
	TokenService    *service.TokenService
	JWTSecret       string
	Config          *config.Config
}

func NewResolver(
	userRepo repository.UserRepository,
	workoutRepo repository.WorkoutRepository,
	exerciseRepo repository.ExerciseRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtSecret string,
	config *config.Config,
) *Resolver {
	return &Resolver{
		UserService:     service.NewUserService(userRepo),
		WorkoutService:  service.NewWorkoutService(workoutRepo),
		ExerciseService: service.NewExerciseService(exerciseRepo),
		TokenService:    service.NewTokenService(refreshTokenRepo),
		JWTSecret:       jwtSecret,
		Config:          config,
	}
}
