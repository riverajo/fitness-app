package seeder

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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

	// 2. Check current version in DB
	metadataColl := db.Collection(MetadataCollection)
	var metadata SystemMetadata
	err = metadataColl.FindOne(ctx, bson.M{"_id": "system_exercises_version"}).Decode(&metadata)

	if err == mongo.ErrNoDocuments {
		// No version found, proceed with seeding
		metadata.Version = 0
	} else if err != nil {
		return fmt.Errorf("failed to fetch system metadata: %w", err)
	}

	if data.Version <= metadata.Version {
		// Already up to date
		return nil
	}

	fmt.Printf("Seeding system exercises (Version %d -> %d)...\n", metadata.Version, data.Version)

	// 3. Seed exercises
	exercisesColl := db.Collection(ExercisesCollection)
	for _, ex := range data.Exercises {
		filter := bson.M{
			"name":   ex.Name,
			"userId": nil, // System exercises have nil UserID
		}

		update := bson.M{
			"$set": bson.M{
				"name":        ex.Name,
				"description": ex.Description,
				"userId":      nil,
				// We might want to add category later to the model, but for now we just use what we have
			},
		}

		opts := options.Update().SetUpsert(true)
		_, err := exercisesColl.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return fmt.Errorf("failed to upsert exercise %s: %w", ex.Name, err)
		}
	}

	// 4. Update version
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
