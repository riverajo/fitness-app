package seeder

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
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

	testDB = client.Database("fitness_db_test")

	code := m.Run()

	if err := mongodbContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	os.Exit(code)
}

func TestSeedSystemExercises(t *testing.T) {
	ctx := context.Background()

	// Helper to create a temp json file
	createTempJSON := func(t *testing.T, data SystemExercisesData) string {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "exercises.json")
		fileContent, err := json.Marshal(data)
		require.NoError(t, err)
		err = os.WriteFile(filePath, fileContent, 0644)
		require.NoError(t, err)
		return filePath
	}

	t.Run("Seeds exercises when DB is empty", func(t *testing.T) {
		require.NoError(t, testDB.Drop(ctx))

		data := SystemExercisesData{
			Version: 1,
			Exercises: []SystemExercise{
				{Name: "Push Up", Description: "Basic push up", Category: "Strength"},
			},
		}
		filePath := createTempJSON(t, data)

		err := SeedSystemExercises(ctx, testDB, filePath)
		require.NoError(t, err)

		// Verify exercise exists
		var result bson.M
		err = testDB.Collection(ExercisesCollection).FindOne(ctx, bson.M{"name": "Push Up"}).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Push Up", result["name"])
		assert.Equal(t, "Basic push up", result["description"])
		assert.Nil(t, result["userId"])

		// Verify version
		var metadata SystemMetadata
		err = testDB.Collection(MetadataCollection).FindOne(ctx, bson.M{"_id": "system_exercises_version"}).Decode(&metadata)
		require.NoError(t, err)
		assert.Equal(t, 1, metadata.Version)
	})

	t.Run("Does not update if version is same or lower", func(t *testing.T) {
		// Setup initial state (Version 1)
		require.NoError(t, testDB.Drop(ctx))
		initialData := SystemExercisesData{
			Version: 1,
			Exercises: []SystemExercise{
				{Name: "Push Up", Description: "Original description", Category: "Strength"},
			},
		}
		err := SeedSystemExercises(ctx, testDB, createTempJSON(t, initialData))
		require.NoError(t, err)

		// Try to seed same version with different data
		newData := SystemExercisesData{
			Version: 1,
			Exercises: []SystemExercise{
				{Name: "Push Up", Description: "Changed description", Category: "Strength"},
			},
		}
		err = SeedSystemExercises(ctx, testDB, createTempJSON(t, newData))
		require.NoError(t, err)

		// Verify description did NOT change
		var result bson.M
		err = testDB.Collection(ExercisesCollection).FindOne(ctx, bson.M{"name": "Push Up"}).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "Original description", result["description"])
	})

	t.Run("Updates exercises when version is higher", func(t *testing.T) {
		// Setup initial state (Version 1)
		require.NoError(t, testDB.Drop(ctx))
		initialData := SystemExercisesData{
			Version: 1,
			Exercises: []SystemExercise{
				{Name: "Push Up", Description: "Original description", Category: "Strength"},
			},
		}
		err := SeedSystemExercises(ctx, testDB, createTempJSON(t, initialData))
		require.NoError(t, err)

		// Seed Version 2
		newData := SystemExercisesData{
			Version: 2,
			Exercises: []SystemExercise{
				{Name: "Push Up", Description: "Updated description", Category: "Strength"},
				{Name: "Pull Up", Description: "New exercise", Category: "Strength"},
			},
		}
		err = SeedSystemExercises(ctx, testDB, createTempJSON(t, newData))
		require.NoError(t, err)

		// Verify updates
		var pushUp bson.M
		err = testDB.Collection(ExercisesCollection).FindOne(ctx, bson.M{"name": "Push Up"}).Decode(&pushUp)
		require.NoError(t, err)
		assert.Equal(t, "Updated description", pushUp["description"])

		var pullUp bson.M
		err = testDB.Collection(ExercisesCollection).FindOne(ctx, bson.M{"name": "Pull Up"}).Decode(&pullUp)
		require.NoError(t, err)
		assert.Equal(t, "New exercise", pullUp["description"])

		// Verify version updated
		var metadata SystemMetadata
		err = testDB.Collection(MetadataCollection).FindOne(ctx, bson.M{"_id": "system_exercises_version"}).Decode(&metadata)
		require.NoError(t, err)
		assert.Equal(t, 2, metadata.Version)
	})

	t.Run("Handles concurrent locking", func(t *testing.T) {
		require.NoError(t, testDB.Drop(ctx))

		// Manually acquire lock
		_, err := testDB.Collection(LocksCollection).InsertOne(ctx, bson.M{
			"_id":       LockID,
			"createdAt": time.Now(),
		})
		require.NoError(t, err)

		// Try to seed (should fail/skip due to lock)
		data := SystemExercisesData{
			Version: 1,
			Exercises: []SystemExercise{
				{Name: "Push Up", Description: "Basic push up", Category: "Strength"},
			},
		}
		err = SeedSystemExercises(ctx, testDB, createTempJSON(t, data))
		require.NoError(t, err) // Should return nil (no error, just skipped)

		// Verify seeding did NOT happen
		count, err := testDB.Collection(ExercisesCollection).CountDocuments(ctx, bson.M{})
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)

		// Release lock
		_, err = testDB.Collection(LocksCollection).DeleteOne(ctx, bson.M{"_id": LockID})
		require.NoError(t, err)

		// Try to seed again (should succeed)
		err = SeedSystemExercises(ctx, testDB, createTempJSON(t, data))
		require.NoError(t, err)

		// Verify seeding happened
		count, err = testDB.Collection(ExercisesCollection).CountDocuments(ctx, bson.M{})
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}
