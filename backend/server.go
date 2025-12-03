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
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/riverajo/fitness-app/backend/graph"
	"github.com/riverajo/fitness-app/backend/internal/config"
	"github.com/riverajo/fitness-app/backend/internal/db"
	"github.com/riverajo/fitness-app/backend/internal/middleware"
	"github.com/riverajo/fitness-app/backend/internal/repository"
	"github.com/riverajo/fitness-app/backend/internal/seeder"
	"github.com/riverajo/fitness-app/backend/internal/spa"
	"github.com/riverajo/fitness-app/backend/telemetry"
)

const defaultPort = "8080"

//go:embed all:public
var publicFS embed.FS

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// 2. Initialize OpenTelemetry
	shutdown, err := telemetry.InitOTel(context.Background(), cfg.AppEnv)
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry", "error", err)
	} else {
		defer func() {
			if err := shutdown(context.Background()); err != nil {
				slog.Error("Failed to shutdown OpenTelemetry", "error", err)
			}
		}()
	}

	// Initialize structured logging with OTel support
	logger := otelslog.NewLogger("fitness-app")
	slog.SetDefault(logger)

	// 3. Connect to Database
	client, err := db.Connect(cfg.MongoURI)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	database := client.Database("fitness_db")

	// Seed System Exercises
	if err := seeder.SeedSystemExercises(context.Background(), database, "data/system_exercises.json"); err != nil {
		slog.Warn("Failed to seed system exercises", "error", err)
	}

	// Initialize Repositories
	userRepo := repository.NewMongoUserRepository(database)
	workoutRepo := repository.NewMongoWorkoutRepository(database)
	exerciseRepo := repository.NewMongoExerciseRepository(database)

	// The Resolver struct is where you inject services like the WorkoutService
	resolver := graph.NewResolver(userRepo, workoutRepo, exerciseRepo, cfg.JWTSecret)

	// 4. GRAPHQL SERVER SETUP
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Wrap the GraphQL server with the necessary middleware chain
	finalHandler := middleware.AuthMiddleware(srv, cfg.JWTSecret)    // 1. Run Auth to validate token and put user ID in context
	finalHandler = middleware.ResponseWriterMiddleware(finalHandler) // 2. Run ResponseWriter injector (needed for setting the cookie)

	// 5. STANDARD GQLGEN CONFIGURATION
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// 6. START SERVER
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
	if cfg.AppEnv == "production" {
		// Production: Serve Embedded SPA + API
		slog.Info("Running in PRODUCTION mode (Single Binary)")

		// Sub-filesystem for "public" folder
		publicFiles, err := fs.Sub(publicFS, "public")
		if err != nil {
			slog.Error("Failed to create sub-filesystem for public", "error", err)
			os.Exit(1)
		}

		// Handle GraphQL API
		http.Handle("/query", otelhttp.NewHandler(finalHandler, "GraphQL"))

		// Handle SPA for everything else
		spaHandler := spa.NewHandler(publicFiles, "index.html")
		http.Handle("/", otelhttp.NewHandler(spaHandler, "SPA"))

	} else {
		// Development: API Only (Playground enabled)
		slog.Info("Running in DEVELOPMENT mode")

		http.Handle("/", otelhttp.NewHandler(playground.Handler("GraphQL playground", "/query"), "Playground"))
		http.Handle("/query", otelhttp.NewHandler(finalHandler, "GraphQL"))
		slog.Info("connect to GraphQL playground", "url", "http://localhost:"+cfg.Port+"/")
	}

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
