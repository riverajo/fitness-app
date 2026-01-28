package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Connect establishes a connection to MongoDB using the URI from environment variables.
func Connect(mongoURI string) (*mongo.Client, error) {
	if mongoURI == "" {
		return nil, fmt.Errorf("mongoURI is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(mongoURI)
	// TODO: Re-enable OpenTelemetry monitoring once the v2 compatible package is available/located.
	// opts.Monitor = otelmongo.NewMonitor()
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the primary to verify the connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	slog.Info("Successfully connected to MongoDB/Cloud Datastore.")
	return client, nil
}
