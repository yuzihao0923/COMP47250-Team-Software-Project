package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var UsersCollection *mongo.Collection

func InitializeMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	db := client.Database("userdb")
	UsersCollection = db.Collection("users")

	// Drop the collection 'users' if it exists
	err = db.Collection("users").Drop(ctx)
	if err != nil && err != mongo.ErrNilDocument {
		return fmt.Errorf("failed to drop existing 'users' collection: %v", err)
	}

	// Recreate the collection and insert initial users
	err = createInitialUsers()
	if err != nil {
		return fmt.Errorf("failed to create initial users: %v", err)
	}

	return nil
}

// connect mongoDB and define UsersCollection
func ConnectMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	UsersCollection = client.Database("userdb").Collection("users")
	return nil
}

func GetUsersCollection() *mongo.Collection {
	return UsersCollection
}

func createInitialUsers() error {
	users := []interface{}{
		bson.M{"username": "broker", "password": "123", "role": "broker"},
		bson.M{"username": "consumer", "password": "123", "role": "consumer"},
		bson.M{"username": "producer", "password": "123", "role": "producer"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := UsersCollection.InsertMany(ctx, users)
	if err != nil {
		return fmt.Errorf("failed to insert initial users: %v", err)
	}

	return nil
}
