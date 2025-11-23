package graph

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/riverajo/fitness-app/backend/graph/model"
	"github.com/riverajo/fitness-app/backend/internal/middleware"
	internalModel "github.com/riverajo/fitness-app/backend/internal/model"
	"github.com/riverajo/fitness-app/backend/internal/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// --- Test Helpers ---

func createTestClient(userRepo *repository.MockUserRepository, workoutRepo *repository.MockWorkoutRepository, exerciseRepo *repository.MockExerciseRepository) *client.Client {
	resolver := NewResolver(userRepo, workoutRepo, exerciseRepo)
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))
	return client.New(srv)
}

// --- Tests ---

func TestRegister(t *testing.T) {
	t.Setenv("JWT_SECRET", "testsecret") // Set environment variable for this test
	userRepo := new(repository.MockUserRepository)
	workoutRepo := new(repository.MockWorkoutRepository)
	exerciseRepo := new(repository.MockExerciseRepository)
	// c := createTestClient(userRepo, workoutRepo, exerciseRepo) // Unused now that we test resolver directly

	input := model.RegisterInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Expect Create to be called
	userRepo.On("Create", mock.Anything, mock.MatchedBy(func(u internalModel.User) bool {
		return u.Email == input.Email && u.PasswordHash != ""
	})).Return(nil)

	// Setup context with ResponseWriter (Need to modify createTestClient or use a custom handler for this test)
	// Since createTestClient hides the server creation, we can't easily inject the middleware there without modifying it.
	// However, gqlgen's client.Post doesn't easily allow injecting context values that the server sees *unless* we use a middleware or a custom transport.
	// A simpler approach for this specific unit test, given the helper, is to instantiate the resolver directly like TestLogin does.
	// This avoids the complexity of the HTTP middleware stack in the test client.

	resolver := NewResolver(userRepo, workoutRepo, exerciseRepo)
	w := httptest.NewRecorder()
	ctx := context.WithValue(context.Background(), "ResponseWriterKey", w)

	payload, err := resolver.Mutation().Register(ctx, input)

	require.NoError(t, err)
	require.True(t, payload.Success)
	require.Equal(t, "Registration successful. You are now logged in.", payload.Message)
	require.Equal(t, input.Email, payload.User.Email)

	// Verify Cookie
	result := w.Result()
	cookies := result.Cookies()
	require.NotEmpty(t, cookies)
	found := false
	for _, c := range cookies {
		if c.Name == middleware.AuthCookieName {
			found = true
			require.NotEmpty(t, c.Value)
			require.True(t, c.HttpOnly)
			require.True(t, c.Secure)
		}
	}
	require.True(t, found, "Auth cookie not found after registration")

	userRepo.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	t.Setenv("JWT_SECRET", "testsecret") // Set environment variable for this test
	userRepo := new(repository.MockUserRepository)
	workoutRepo := new(repository.MockWorkoutRepository)
	exerciseRepo := new(repository.MockExerciseRepository)
	resolver := NewResolver(userRepo, workoutRepo, exerciseRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &internalModel.User{
		ID:           "user123",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}

	userRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)

	// Setup context with ResponseWriter
	w := httptest.NewRecorder()
	ctx := context.WithValue(context.Background(), "ResponseWriterKey", w)

	input := model.LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	payload, err := resolver.Mutation().Login(ctx, input)

	require.NoError(t, err)
	require.True(t, payload.Success)
	require.Equal(t, "Login successful. Token set in cookie.", payload.Message)
	require.Equal(t, user.ID, payload.User.ID)

	// Verify Cookie
	result := w.Result()
	cookies := result.Cookies()
	require.NotEmpty(t, cookies)
	found := false
	for _, c := range cookies {
		if c.Name == middleware.AuthCookieName {
			found = true
			require.NotEmpty(t, c.Value)
			require.True(t, c.HttpOnly)
			require.True(t, c.Secure)
		}
	}
	require.True(t, found, "Auth cookie not found")
	userRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	userRepo := new(repository.MockUserRepository)
	workoutRepo := new(repository.MockWorkoutRepository)
	exerciseRepo := new(repository.MockExerciseRepository)
	resolver := NewResolver(userRepo, workoutRepo, exerciseRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &internalModel.User{
		ID:           "user123",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}

	// Setup context with UserID
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "user123")

	// Mock GetUserByID (called by UpdateUser to verify current password)
	userRepo.On("FindByID", mock.Anything, "user123").Return(user, nil)

	// Mock Update
	userRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *internalModel.User) bool {
		return u.ID == "user123" && u.PreferredUnit == "kg"
	})).Return(nil)

	currentPwd := "password123"
	newUnit := "kg"
	input := model.UpdateUserInput{
		CurrentPassword: currentPwd,
		PreferredUnit:   &newUnit,
	}

	payload, err := resolver.Mutation().UpdateUser(ctx, input)

	require.NoError(t, err)
	require.True(t, payload.Success)
	require.Equal(t, "kg", payload.User.PreferredUnit)
	userRepo.AssertExpectations(t)
}

func TestMe(t *testing.T) {
	userRepo := new(repository.MockUserRepository)
	workoutRepo := new(repository.MockWorkoutRepository)
	exerciseRepo := new(repository.MockExerciseRepository)

	user := &internalModel.User{
		ID:    "user123",
		Email: "test@example.com",
	}

	userRepo.On("FindByID", mock.Anything, "user123").Return(user, nil)

	resolver := NewResolver(userRepo, workoutRepo, exerciseRepo)
	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "user123")

	me, err := resolver.Query().Me(ctx)

	require.NoError(t, err)
	require.Equal(t, "user123", me.ID)
	require.Equal(t, "test@example.com", me.Email)
	userRepo.AssertExpectations(t)
}

func TestCreateWorkoutLog(t *testing.T) {
	userRepo := new(repository.MockUserRepository)
	workoutRepo := new(repository.MockWorkoutRepository)
	exerciseRepo := new(repository.MockExerciseRepository)
	resolver := NewResolver(userRepo, workoutRepo, exerciseRepo)

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "user123")

	input := model.CreateWorkoutLogInput{
		Name: "Morning Workout",
		ExerciseLogs: []*model.ExerciseLogInput{
			{
				UniqueExerciseID: "ex1",
				Sets: []*model.SetInput{
					{Reps: 10, Weight: 100},
				},
			},
		},
	}

	expectedLog := &internalModel.WorkoutLog{
		ID:     "log123",
		UserID: "user123",
		Name:   "Morning Workout",
	}

	workoutRepo.On("Create", mock.Anything, mock.MatchedBy(func(l internalModel.WorkoutLog) bool {
		return l.UserID == "user123" && l.Name == "Morning Workout" && len(l.ExerciseLogs) == 1
	})).Return(expectedLog, nil)

	log, err := resolver.Mutation().CreateWorkoutLog(ctx, input)

	require.NoError(t, err)
	require.Equal(t, "log123", log.ID)
	workoutRepo.AssertExpectations(t)
}

func TestGetWorkoutLog(t *testing.T) {
	userRepo := new(repository.MockUserRepository)
	workoutRepo := new(repository.MockWorkoutRepository)
	exerciseRepo := new(repository.MockExerciseRepository)
	resolver := NewResolver(userRepo, workoutRepo, exerciseRepo)

	expectedLog := &internalModel.WorkoutLog{
		ID:   "log123",
		Name: "Test Log",
	}

	workoutRepo.On("GetByID", mock.Anything, "log123").Return(expectedLog, nil)

	log, err := resolver.Query().GetWorkoutLog(context.Background(), "log123")

	require.NoError(t, err)
	require.Equal(t, "log123", log.ID)
	workoutRepo.AssertExpectations(t)
}

func TestListWorkoutLogs(t *testing.T) {
	userRepo := new(repository.MockUserRepository)
	workoutRepo := new(repository.MockWorkoutRepository)
	exerciseRepo := new(repository.MockExerciseRepository)
	resolver := NewResolver(userRepo, workoutRepo, exerciseRepo)

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "user123")

	expectedLogs := []*internalModel.WorkoutLog{
		{ID: "log1", Name: "Log 1"},
		{ID: "log2", Name: "Log 2"},
	}

	limit := 10
	offset := 0
	workoutRepo.On("ListByUser", mock.Anything, "user123", limit, offset).Return(expectedLogs, nil)

	logs, err := resolver.Query().ListWorkoutLogs(ctx, int32Ptr(int32(limit)), int32Ptr(int32(offset)))

	require.NoError(t, err)
	require.Len(t, logs, 2)
	workoutRepo.AssertExpectations(t)
}

func TestLogoutMutation(t *testing.T) {
	// 1. Setup
	// Logout doesn't use any services, so we can initialize Resolver with nil services
	resolver := &Resolver{}

	// Create a ResponseRecorder to capture the cookie
	w := httptest.NewRecorder()

	// Create a context with the ResponseWriter injected (using the key from middleware)
	ctx := context.WithValue(context.Background(), "ResponseWriterKey", w)

	// 2. Execute
	payload, err := resolver.Mutation().Logout(ctx)

	// 3. Verify
	if err != nil {
		t.Fatalf("Logout returned error: %v", err)
	}

	if !payload.Success {
		t.Error("Expected success to be true")
	}

	if payload.Message != "Logged out successfully." {
		t.Errorf("Unexpected message: %s", payload.Message)
	}

	// Check the cookie
	result := w.Result()
	cookies := result.Cookies()

	found := false
	for _, cookie := range cookies {
		if cookie.Name == middleware.AuthCookieName {
			found = true
			if cookie.MaxAge != -1 {
				t.Errorf("Expected MaxAge to be -1, got %d", cookie.MaxAge)
			}
			if cookie.Value != "" {
				t.Errorf("Expected cookie value to be empty, got %s", cookie.Value)
			}
			break
		}
	}

	if !found {
		t.Error("auth_token cookie was not set")
	}
}

func TestCreateUniqueExercise(t *testing.T) {
	userRepo := new(repository.MockUserRepository)
	workoutRepo := new(repository.MockWorkoutRepository)
	exerciseRepo := new(repository.MockExerciseRepository)
	resolver := NewResolver(userRepo, workoutRepo, exerciseRepo)

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "user123")

	input := model.CreateUniqueExerciseInput{
		Name: "My Custom Bench",
	}

	exerciseRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *internalModel.UniqueExercise) bool {
		return e.Name == "My Custom Bench" && *e.UserID == "user123"
	})).Return(nil).Run(func(args mock.Arguments) {
		e := args.Get(1).(*internalModel.UniqueExercise)
		e.ID = "ex123"
	})

	ex, err := resolver.Mutation().CreateUniqueExercise(ctx, input)

	require.NoError(t, err)
	require.Equal(t, "My Custom Bench", ex.Name)
	require.Equal(t, "user123", *ex.UserID)
	exerciseRepo.AssertExpectations(t)
}

func TestUniqueExercises(t *testing.T) {
	userRepo := new(repository.MockUserRepository)
	workoutRepo := new(repository.MockWorkoutRepository)
	exerciseRepo := new(repository.MockExerciseRepository)
	resolver := NewResolver(userRepo, workoutRepo, exerciseRepo)

	ctx := context.WithValue(context.Background(), middleware.UserIDKey, "user123")

	expectedExercises := []*internalModel.UniqueExercise{
		{ID: "ex1", Name: "Bench Press", UserID: nil},
		{ID: "ex2", Name: "My Bench", UserID: stringPtr("user123")},
	}

	// Expect Search with default limit 50 and offset 0
	exerciseRepo.On("Search", mock.Anything, stringPtr("user123"), "Bench", 50, 0).Return(expectedExercises, nil)

	query := "Bench"
	limit := int32(50)
	offset := int32(0)
	exercises, err := resolver.Query().UniqueExercises(ctx, &query, &limit, &offset)

	require.NoError(t, err)
	require.Len(t, exercises, 2)
	exerciseRepo.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}
