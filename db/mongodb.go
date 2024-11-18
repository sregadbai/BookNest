package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// ConnectMongoDB initializes the MongoDB client
func ConnectMongoDB(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Verify the connection
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB!")
}

// GetCollection retrieves a collection from the database
func GetCollection(collectionName string) *mongo.Collection {
	if Client == nil {
		log.Fatal("MongoDB client is not initialized. Ensure ConnectMongoDB is called before accessing collections.")
	}
	return Client.Database("booknest").Collection(collectionName)
}
