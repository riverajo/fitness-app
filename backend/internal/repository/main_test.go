package repository

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func getMyNetworkName(ctx context.Context, containerID string) string {
	if containerID == "" {
		return ""
	}

	// Create a provider to talk to the socket
	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		return ""
	}
	defer func() {
		if err := provider.Close(); err != nil {
			log.Printf("failed to close provider: %s", err)
		}
	}()

	// Inspect "this" container
	client, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return ""
	}

	inspect, err := client.ContainerInspect(ctx, containerID)
	if err != nil {
		return ""
	}

	// Grab the first network name attached to this container
	for netName := range inspect.NetworkSettings.Networks {
		return netName
	}
	return ""
}

func getMongoURI(ctx context.Context, container *mongodb.MongoDBContainer, networkName string) string {
	// If we are in a shared network (CI), use the internal alias and port
	if networkName != "" && networkName != "bridge" {
		return "mongodb://mongodb_repo:27017/fitness_db?authSource=admin"
	}

	// Otherwise (Local Dev), use the official helper
	uri, err := container.ConnectionString(ctx)
	if err != nil {
		return "mongodb://localhost:27017/fitness_db?authSource=admin"
	}
	return uri
}

var testDB *mongo.Database

func TestMain(m *testing.M) {
	ctx := context.Background()

	myID := os.Getenv("MY_CONTAINER_ID")
	networkName := getMyNetworkName(ctx, myID)

	var opts []testcontainers.ContainerCustomizer

	if networkName != "" {
		// Use the network package's helper to join by name
		opts = append(opts, network.WithNetworkName([]string{"mongodb_repo"}, networkName),
			testcontainers.WithWaitStrategy(
				wait.ForAll(
					wait.ForListeningPort("27017/tcp").SkipExternalCheck(),
				).WithDeadline(2*time.Minute), // Changed from WithStartupTimeout
			),
		)
	}
	mongodbContainer, err := mongodb.Run(ctx, "mongo:6", opts...)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	endpoint := getMongoURI(ctx, mongodbContainer, networkName)
	client, err := mongo.Connect(options.Client().ApplyURI(endpoint))
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
