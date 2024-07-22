package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *MongoDB

type MongoDB struct {
	Client          *mongo.Client
	UsersCollection *mongo.Collection
}

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Role     string `bson:"role"`
}

func NewMongoDB(uri, dbName, collectionName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	dbClient = &MongoDB{
		Client:          client,
		UsersCollection: collection,
	}

	if err := dbClient.InitializeMongoDB(ctx); err != nil {
		return nil, err
	}

	return dbClient, nil
}

func (db *MongoDB) InitializeMongoDB(ctx context.Context) error {
	if err := db.UsersCollection.Drop(ctx); err != nil && err != mongo.ErrNilDocument {
		return fmt.Errorf("failed to drop existing 'users' collection: %v", err)
	}
	return db.createInitialUsers()
}

func (db *MongoDB) createInitialUsers() error {
	users := []interface{}{
		bson.M{"username": "b1", "password": "123", "role": "broker"},
		bson.M{"username": "c1", "password": "123", "role": "consumer"},
		bson.M{"username": "c2", "password": "123", "role": "consumer"},
		bson.M{"username": "c3", "password": "123", "role": "consumer"},
		bson.M{"username": "p1", "password": "123", "role": "producer"},
		bson.M{"username": "p2", "password": "123", "role": "producer"},
		bson.M{"username": "p3", "password": "123", "role": "producer"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.UsersCollection.InsertMany(ctx, users)
	if err != nil {
		return fmt.Errorf("failed to insert initial users: %v", err)
	}

	return nil
}

func (db *MongoDB) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := db.UsersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *MongoDB) Close(ctx context.Context) error {
	return db.Client.Disconnect(ctx)
}

func ConnectMongoDB(uri, dbName, collectionName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	dbClient = &MongoDB{
		Client:          client,
		UsersCollection: collection,
	}

	return dbClient, nil
}

func GetDBClient() *MongoDB {
	return dbClient
}
