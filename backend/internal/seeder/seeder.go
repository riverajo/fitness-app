package seeder

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SystemExercise struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type SystemExercisesData struct {
	Version   int              `json:"version"`
	Exercises []SystemExercise `json:"exercises"`
}

type SystemMetadata struct {
	ID      string `bson:"_id,omitempty"`
	Version int    `bson:"version"`
}

const MetadataCollection = "system_metadata"
const ExercisesCollection = "unique_exercises"
const LocksCollection = "system_locks"
const LockID = "seeder_lock"

// ensureLockIndex creates a TTL index on the locks collection to handle crashes
func ensureLockIndex(ctx context.Context, db *mongo.Database) error {
	locksColl := db.Collection(LocksCollection)
	_, err := locksColl.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"createdAt": 1},
		Options: options.Index().SetExpireAfterSeconds(60), // Lock expires after 60 seconds
	})
	return err
}

// SeedSystemExercises loads system exercises from a JSON file if the version is newer than what's in the DB.
func SeedSystemExercises(ctx context.Context, db *mongo.Database, filePath string) error {
	// 1. Read JSON file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read system exercises file: %w", err)
	}

	var data SystemExercisesData
	if err := json.Unmarshal(fileContent, &data); err != nil {
		return fmt.Errorf("failed to parse system exercises json: %w", err)
	}

	// 2. Check current version in DB (Optimization: check before locking)
	metadataColl := db.Collection(MetadataCollection)
	var metadata SystemMetadata
	err = metadataColl.FindOne(ctx, bson.M{"_id": "system_exercises_version"}).Decode(&metadata)

	if err == mongo.ErrNoDocuments {
		metadata.Version = 0
	} else if err != nil {
		return fmt.Errorf("failed to fetch system metadata: %w", err)
	}

	if data.Version <= metadata.Version {
		return nil
	}

	// 3. Acquire Lock
	if err := ensureLockIndex(ctx, db); err != nil {
		return fmt.Errorf("failed to ensure lock index: %w", err)
	}

	locksColl := db.Collection(LocksCollection)
	_, err = locksColl.InsertOne(ctx, bson.M{
		"_id":       LockID,
		"createdAt": time.Now(),
	})

	if mongo.IsDuplicateKeyError(err) {
		fmt.Println("Seeding in progress by another node. Skipping.")
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	// Ensure lock is released
	defer func() {
		_, _ = locksColl.DeleteOne(context.Background(), bson.M{"_id": LockID})
	}()

	// Re-check version inside lock to be safe (double-checked locking)
	err = metadataColl.FindOne(ctx, bson.M{"_id": "system_exercises_version"}).Decode(&metadata)
	if err == mongo.ErrNoDocuments {
		metadata.Version = 0
	} else if err != nil {
		return fmt.Errorf("failed to fetch system metadata: %w", err)
	}

	if data.Version <= metadata.Version {
		return nil
	}

	fmt.Printf("Seeding system exercises (Version %d -> %d)...\n", metadata.Version, data.Version)

	// 4. Seed exercises
	exercisesColl := db.Collection(ExercisesCollection)
	for _, ex := range data.Exercises {
		filter := bson.M{
			"name":   ex.Name,
			"userId": nil,
		}

		update := bson.M{
			"$set": bson.M{
				"name":        ex.Name,
				"description": ex.Description,
				"userId":      nil,
			},
		}

		opts := options.Update().SetUpsert(true)
		_, err := exercisesColl.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return fmt.Errorf("failed to upsert exercise %s: %w", ex.Name, err)
		}
	}

	// 5. Update version
	_, err = metadataColl.UpdateOne(
		ctx,
		bson.M{"_id": "system_exercises_version"},
		bson.M{"$set": bson.M{"version": data.Version}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to update system version: %w", err)
	}

	fmt.Println("System exercises seeded successfully.")
	return nil
}
