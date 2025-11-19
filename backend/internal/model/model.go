package model

import (
	"time"
)

// --- INPUT TYPES ---

// LoginInput matches the GraphQL input type for user login.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Maps to the GraphQL 'ExerciseLogInput'
type ExerciseLogInput struct {
	UniqueExerciseID string      `json:"uniqueExerciseId"`
	Sets             []*SetInput `json:"sets"`
	Notes            *string     `json:"notes"`
}

// Maps to the GraphQL 'SetInput'
type SetInput struct {
	Reps      int32   `json:"reps"`
	Weight    float64 `json:"weight"`
	Rpe       *int32  `json:"rpe"` // Pointer for optional field
	ToFailure *bool   `json:"toFailure"`
	Order     int32   `json:"order"`
}

// Maps to the GraphQL 'CreateWorkoutLogInput'
type CreateWorkoutLogInput struct {
	Name         string              `json:"name"`
	StartTime    time.Time           `json:"startTime"`
	EndTime      time.Time           `json:"endTime"`
	ExerciseLogs []*ExerciseLogInput `json:"exerciseLogs"`
	LocationName *string             `json:"locationName"`
	GeneralNotes *string             `json:"generalNotes"`
}

// --- CORE DATABASE MODEL (Used for Output) ---

// Maps to the GraphQL 'WorkoutLog' and MongoDB storage
type WorkoutLog struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	// ... other fields matching the GraphQL schema...
}

const (
	WeightUnitKilograms = "KGS"
	WeightUnitPounds    = "LBS"
)
