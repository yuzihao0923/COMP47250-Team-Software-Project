package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *MongoDB // There wont be concurrent issue for dbClient, cause ConnectMongoDB func runs on differet main func, which means they are in different stacks

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //10秒时间来连接数据库
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	return &MongoDB{
		Client:          client,
		UsersCollection: collection,
	}, nil
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

func ConnectMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	db := client.Database("yourdbname")
	collection := db.Collection("yourcollectionname")

	dbClient = &MongoDB{
		Client:          client,
		UsersCollection: collection,
	}

	return nil
}

func GetDBClient() *MongoDB {
	return dbClient
}

// var client *mongo.Client
// var UsersCollection *mongo.Collection

// func InitializeMongoDB() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	var err error
// 	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to MongoDB: %v", err)
// 	}

// 	db := client.Database("userdb")
// 	UsersCollection = db.Collection("users")

// 	// Drop the collection 'users' if it exists
// 	err = db.Collection("users").Drop(ctx)
// 	if err != nil && err != mongo.ErrNilDocument {
// 		return fmt.Errorf("failed to drop existing 'users' collection: %v", err)
// 	}

// 	// Recreate the collection and insert initial users
// 	err = createInitialUsers()
// 	if err != nil {
// 		return fmt.Errorf("failed to create initial users: %v", err)
// 	}

// 	return nil
// }

// connect mongoDB and define UsersCollection
// func ConnectMongoDB() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	var err error
// 	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to MongoDB: %v", err)
// 	}

// 	UsersCollection = client.Database("userdb").Collection("users")
// 	return nil
// }

// func GetUsersCollection() *mongo.Collection {
// 	return UsersCollection
// }

// func createInitialUsers() error {
// 	users := []interface{}{
// 		bson.M{"username": "b1", "password": "123", "role": "broker"},
// 		bson.M{"username": "c1", "password": "123", "role": "consumer"},
// 		bson.M{"username": "c2", "password": "123", "role": "consumer"},
// 		bson.M{"username": "c3", "password": "123", "role": "consumer"},
// 		bson.M{"username": "p1", "password": "123", "role": "producer"},
// 		bson.M{"username": "p2", "password": "123", "role": "producer"},
// 		bson.M{"username": "p3", "password": "123", "role": "producer"},
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	_, err := UsersCollection.InsertMany(ctx, users)
// 	if err != nil {
// 		return fmt.Errorf("failed to insert initial users: %v", err)
// 	}

// 	return nil
// }
