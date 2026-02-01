package main

import (
	"context"
	"embed"

	"io/fs"
	"log/slog"
	"net/http"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/grafana/pyroscope-go"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/vektah/gqlparser/v2/ast"

	"github.com/riverajo/fitness-app/backend/graph"
	"github.com/riverajo/fitness-app/backend/internal/api"
	"github.com/riverajo/fitness-app/backend/internal/config"
	"github.com/riverajo/fitness-app/backend/internal/db"
	"github.com/riverajo/fitness-app/backend/internal/middleware"
	"github.com/riverajo/fitness-app/backend/internal/repository"
	"github.com/riverajo/fitness-app/backend/internal/seeder"
	"github.com/riverajo/fitness-app/backend/internal/spa"
	"github.com/riverajo/fitness-app/backend/telemetry"
)

//go:embed all:public
var publicFS embed.FS

//go:embed data/system_exercises.json
var systemExercisesData []byte

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// 2. Initialize OpenTelemetry
	shutdown, err := telemetry.InitOTel(context.Background(), cfg.AppEnv, cfg.CI)
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry", "error", err)
	} else {
		defer func() {
			if err := shutdown(context.Background()); err != nil {
				slog.Error("Failed to shutdown OpenTelemetry", "error", err)
			}
		}()
	}

	// 3. Initialize Pyroscope (if configured)
	if cfg.PyroscopeURL != "" {
		slog.Info("Initializing Pyroscope", "url", cfg.PyroscopeURL, "app_name", cfg.PyroscopeAppName)
		_, err := pyroscope.Start(pyroscope.Config{
			ApplicationName: cfg.PyroscopeAppName,
			ServerAddress:   cfg.PyroscopeURL,
			Logger:          nil, // Use default logger or nil to avoid noise
			ProfileTypes: []pyroscope.ProfileType{
				pyroscope.ProfileCPU,
				pyroscope.ProfileAllocObjects,
				pyroscope.ProfileAllocSpace,
				pyroscope.ProfileInuseObjects,
				pyroscope.ProfileInuseSpace,
			},
		})
		if err != nil {
			slog.Error("Failed to start Pyroscope", "error", err)
		}
	} else {
		slog.Info("Pyroscope disabled (PYROSCOPE_URL not set)")
	}

	// Initialize structured logging with OTel support
	logger := otelslog.NewLogger("fitness-app")
	slog.SetDefault(logger)

	// 3. Connect to Database
	client, err := db.Connect(cfg.MongoURI)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		// Ensure logs are flushed before exiting
		if shutdown != nil {
			_ = shutdown(context.Background())
		}
		os.Exit(1)
	}

	database := client.Database("fitness_db")

	// Seed System Exercises
	if err := seeder.SeedSystemExercises(context.Background(), database, systemExercisesData); err != nil {
		slog.Warn("Failed to seed system exercises", "error", err)
	}

	// Initialize Repositories
	userRepo := repository.NewMongoUserRepository(database)
	workoutRepo := repository.NewMongoWorkoutRepository(database)
	exerciseRepo := repository.NewMongoExerciseRepository(database)
	refreshTokenRepo := repository.NewMongoRefreshTokenRepository(database)

	// The Resolver struct is where you inject services like the WorkoutService
	resolver := graph.NewResolver(userRepo, workoutRepo, exerciseRepo, refreshTokenRepo, cfg.JWTSecret, cfg)

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

	// Handle Faro Collection Proxy
	http.Handle("/faro/collect", otelhttp.NewHandler(middleware.FaroProxy(cfg.FaroURL, cfg.CI), "FaroProxy"))

	// 7. AUTH HANDLERS
	secureCookie := cfg.AppEnv == "production" && !cfg.CI
	authHandler := api.NewAuthHandler(resolver.TokenService, resolver.UserService, cfg.JWTSecret, secureCookie)
	http.HandleFunc("/auth/refresh", authHandler.Refresh)

	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := client.Ping(r.Context(), nil); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("db not ready"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	finalHandler = middleware.RecoveryMiddleware(finalHandler)
	finalHandler = middleware.LoggingMiddleware(finalHandler)

	// Determine environment
	if cfg.AppEnv == "production" {
		// Production: Serve Embedded SPA + API
		slog.Info("Running in PRODUCTION mode (Single Binary)")

		// Sub-filesystem for "public" folder
		publicFiles, err := fs.Sub(publicFS, "public")
		if err != nil {
			slog.Error("Failed to create sub-filesystem for public", "error", err)
			if shutdown != nil {
				_ = shutdown(context.Background())
			}
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

	// 6. START SERVER
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: nil, // Use DefaultServeMux
	}

	// Run server in a goroutine so it doesn't block
	go func() {
		slog.Info("Server starting", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			if shutdown != nil {
				_ = shutdown(context.Background())
			}
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown signal received")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	// Clean up resources
	if shutdown != nil {
		if err := shutdown(context.Background()); err != nil {
			slog.Error("Failed to shutdown OpenTelemetry", "error", err)
		}
	}

	if err := client.Disconnect(context.Background()); err != nil {
		slog.Error("Failed to disconnect from database", "error", err)
	}

	slog.Info("Shutdown complete")
}
