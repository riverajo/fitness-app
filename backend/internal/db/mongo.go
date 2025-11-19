package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient holds the connection pool
var MongoClient *mongo.Client

// Connect establishes a connection to MongoDB using the URI from environment variables.
func Connect() {
	// The MONGO_URI is set in the docker-compose.yml (or GCP environment)
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI environment variable not set.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the primary to verify the connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB/Cloud Datastore.")
	MongoClient = client
}

// GetCollection returns a specific collection from the database
func GetCollection(collectionName string) *mongo.Collection {
	// We hardcode the database name, which should be consistent with the URI
	return MongoClient.Database("fitness_db").Collection(collectionName)
}
