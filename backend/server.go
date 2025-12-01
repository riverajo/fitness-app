package main

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/riverajo/fitness-app/backend/graph"
	"github.com/riverajo/fitness-app/backend/internal/db"
	"github.com/riverajo/fitness-app/backend/internal/middleware"
	"github.com/riverajo/fitness-app/backend/internal/repository"
	"github.com/riverajo/fitness-app/backend/internal/seeder"
	"github.com/riverajo/fitness-app/backend/internal/spa"
)

const defaultPort = "8080"

//go:embed public
var publicFS embed.FS

func main() {

	// 1. 💡 DATABASE CONNECTION: Establish connection to MongoDB/Cloud Datastore
	// This function will look for MONGO_URI and fatal if it fails to connect.
	// Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	client, err := db.Connect()
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	database := client.Database("fitness_db")

	// Seed System Exercises
	if err := seeder.SeedSystemExercises(context.Background(), database, "data/system_exercises.json"); err != nil {
		slog.Warn("Failed to seed system exercises", "error", err)
		// We don't fatal here because the server should still run even if seeding fails
	}

	// Initialize Repositories
	userRepo := repository.NewMongoUserRepository(database)
	workoutRepo := repository.NewMongoWorkoutRepository(database)
	exerciseRepo := repository.NewMongoExerciseRepository(database)

	// The Resolver struct is where you inject services like the WorkoutService
	resolver := graph.NewResolver(userRepo, workoutRepo, exerciseRepo)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// 3. GRAPHQL SERVER SETUP (Uses the generated code and our custom resolver)
	// The original line: srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	// Is now updated to use our prepared resolver:
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	// Wrap the GraphQL server with the necessary middleware chain
	// The outer-most handler (last one applied) runs first.
	finalHandler := middleware.AuthMiddleware(srv)                   // 1. Run Auth to validate token and put user ID in context
	finalHandler = middleware.ResponseWriterMiddleware(finalHandler) // 2. Run ResponseWriter injector (needed for setting the cookie)

	// 4. STANDARD GQLGEN CONFIGURATION (Keep these from the generated file)
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// 5. START SERVER
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := client.Ping(r.Context(), nil); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("db not ready"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	// Determine environment
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		// Production: Serve Embedded SPA + API
		slog.Info("Running in PRODUCTION mode (Single Binary)")

		// Sub-filesystem for "public" folder
		publicFiles, err := fs.Sub(publicFS, "public")
		if err != nil {
			slog.Error("Failed to create sub-filesystem for public", "error", err)
			os.Exit(1)
		}

		// Handle GraphQL API
		http.Handle("/query", finalHandler)

		// Handle SPA for everything else
		spaHandler := spa.NewHandler(publicFiles, "index.html")
		http.Handle("/", spaHandler)

	} else {
		// Development: API Only (Playground enabled)
		slog.Info("Running in DEVELOPMENT mode")

		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", finalHandler)
	}

	slog.Info("connect to GraphQL playground", "url", "http://localhost:"+port+"/")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
