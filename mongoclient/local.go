package mongoclient

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func ConnectLocal() *mongo.Client {
	username := os.Getenv("MONGODB_USERNAME")
	if username == "" {
		username = "root"
	}
	password := os.Getenv("MONGODB_PASSWORD")
	if password == "" {
		password = "12345"
	}
	port := os.Getenv("MONGODB_LOCAL_PORT")
	if port == "" {
		port = "27018"
	}

	url := fmt.Sprintf("mongodb://%s:%s@localhost:%v/admin?directConnection=true&serverSelectionTimeoutMS=2000&authSource=admin", username, password, port)
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}
