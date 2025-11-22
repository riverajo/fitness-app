package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testDB *mongo.Database

func TestMain(m *testing.M) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		log.Fatalf("failed to connect to mongo: %s", err)
	}

	// Set the global testDB
	testDB = client.Database("fitness_db")

	code := m.Run()

	if err := mongodbContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	os.Exit(code)
}

// Helper to clean up collection after each test
func cleanupCollection(t *testing.T, collectionName string) {
	t.Helper()
	if testDB == nil {
		t.Fatal("testDB is nil")
	}
	err := testDB.Collection(collectionName).Drop(context.Background())
	if err != nil {
		t.Fatalf("failed to drop collection %s: %v", collectionName, err)
	}
}
