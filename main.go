package main

import (
	"context"
	"fmt"
	"github.com/ArnobKumarSaha/mongo/database"
	"github.com/ArnobKumarSaha/mongo/latency"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func connectToMongo() *mongo.Client {
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")

	url := fmt.Sprintf("mongodb://%s:%s@localhost:27018/admin?directConnection=true&serverSelectionTimeoutMS=2000&authSource=admin", username, password)
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected")
	return client
}

func main() {
	client := connectToMongo()
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	collMap := database.ListCollectionsForAllDatabases(client)
	for db, collections := range collMap {
		for _, coll := range collections {
			latency.CalculateLatency(client.Database(db), coll)
		}
	}
	latency.PrintLatencyStats()
}
