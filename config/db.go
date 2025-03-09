package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}

	// Ping the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Failed to ping MongoDB:", err)
	}

	fmt.Println("✅ Connected to MongoDB!")

	DB = client.Database("go_fiber")

	// Define which collections and fields require unique indexes
	uniqueFields := map[string][]string{
		"users":      {"email"},
		"products":   {"name"},
		"categories": {"category_name"},
		"customers":  {"email"},
		"orders":     {"order_number"},
		// Add more collections and fields as needed
	}

	// Automatically create unique indexes based on the map
	err = createUniqueIndexesForCollections(uniqueFields)
	if err != nil {
		log.Fatal("Error creating unique index:", err)
	}
}

func createUniqueIndexesForCollections(uniqueFields map[string][]string) error {
	// Loop through the map to get collection names and unique fields
	for collectionName, fields := range uniqueFields {
		// Get the collection reference
		collection := DB.Collection(collectionName)

		for _, field := range fields {
			// Create a unique index for each field
			indexModel := mongo.IndexModel{
				Keys:    bson.M{field: 1}, // Create index on the field
				Options: options.Index().SetUnique(true),
			}

			// Apply the index to the collection
			_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				log.Printf("Could not create unique index for collection %s and field %s: %v", collectionName, field, err)
				continue // Move on to the next field/collection
			}
			log.Printf("✅ Unique index on %s field created for collection %s", field, collectionName)
		}
	}
	return nil
}
