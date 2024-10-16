package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

// GetMongoClient returns a singleton MongoDB client instance
func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		// Load environment variables
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found. Using environment variables.")
		}

		mongoURI := os.Getenv("MONGODB_URI")
		if mongoURI == "" {
			clientInstanceError = fmt.Errorf("MONGODB_URI is not set in environment variables")
			return
		}

		clientOptions := options.Client().ApplyURI(mongoURI)
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
			return
		}

		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
			return
		}

		log.Println("Connected to MongoDB Atlas!")
		clientInstance = client
	})

	return clientInstance, clientInstanceError
}
