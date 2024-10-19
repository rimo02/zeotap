package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var WeatherClient *mongo.Client
var UserClient *mongo.Client

func ConnectDB(uri string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB: ", err)
	}

	fmt.Printf("Connected to MongoDB at %s\n", uri)
	return client
}

func InitializeConnections() {
	WeatherClient = ConnectDB("mongodb://localhost:27017") // Connect to the weather trends database
	UserClient = ConnectDB("mongodb://localhost:27017")    // Connect to the users database
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	db := client.Database("weather")
	return getOrCreateCollection(db, collectionName)
}

func GetUserCollection(client *mongo.Client, city string) *mongo.Collection {
	db := client.Database("users")
	return getOrCreateCollection(db, city)
}

func getOrCreateCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	exist, err := db.ListCollectionNames(context.TODO(), bson.M{"name": collectionName})
	if err != nil {
		fmt.Printf("Error listing collections: %v\n", err)
		return nil
	}

	if len(exist) == 0 {
		err := db.CreateCollection(context.TODO(), collectionName)
		if err != nil {
			fmt.Printf("Error creating collection: %v\n", err)
			return nil
		}
	}

	return db.Collection(collectionName)
}
