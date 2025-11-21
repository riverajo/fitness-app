package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/riverajo/fitness-app/backend/graph"
	"github.com/riverajo/fitness-app/backend/internal/db"
	"github.com/riverajo/fitness-app/backend/internal/middleware"
	"github.com/riverajo/fitness-app/backend/internal/repository"
)

const defaultPort = "8080"

func main() {

	// 1. 💡 DATABASE CONNECTION: Establish connection to MongoDB/Cloud Datastore
	// This function will look for MONGO_URI and fatal if it fails to connect.
	db.Connect()

	// Initialize Repositories
	userRepo := repository.NewMongoUserRepository()
	workoutRepo := repository.NewMongoWorkoutRepository()

	// The Resolver struct is where you inject services like the WorkoutService
	resolver := graph.NewResolver(userRepo, workoutRepo)

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

	allowedOrigins := []string{
		// 🚨 REPLACE THIS WITH YOUR PRODUCTION SVELTEKIT DOMAIN
		"https://app.your-domain.com",

		// For local development
		"http://localhost:3000",
		"http://localhost:5173",
		"http://localhost:4173", // Playwright preview port
		"http://localhost:8080",
	}

	// Create the CORS handler configuration
	c := cors.New(cors.Options{
		// Which origins are allowed to make requests
		AllowedOrigins: allowedOrigins,

		// Methods required by GraphQL (POST) and preflight checks (OPTIONS)
		AllowedMethods: []string{http.MethodPost, http.MethodOptions},

		// Headers required by GraphQL (Content-Type) and Auth (Authorization/Cookie)
		AllowedHeaders: []string{"Content-Type", "Authorization", "Cookie"},

		// CRITICAL: Allows the browser to send cookies/auth headers
		AllowCredentials: true,

		// Optional: Preflight cache time
		MaxAge: 300,
	})

	// 5. START SERVER
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", c.Handler(finalHandler))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
