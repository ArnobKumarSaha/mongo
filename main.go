package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/klog/v2"
	"log"
	"os"
	"sort"
)

func connectToMongo() *mongo.Client {
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")

	url := fmt.Sprintf("mongodb://%s:%s@localhost:27017/admin?directConnection=true&serverSelectionTimeoutMS=2000&authSource=admin", username, password)
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected")
	return client
}

type read struct {
	latency, ops int64
	db, coll     string
}

var reads []read

func main() {
	client := connectToMongo()
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	//db := client.Database("test")
	//test(db, "kubedb")

	dbs, err := client.ListDatabases(context.TODO(), bson.D{})
	if err != nil {
		klog.Errorf("err = %s \n", err.Error())
		return
	}
	for _, d := range dbs.Databases {
		db := client.Database(d.Name)
		names, err := db.ListCollectionNames(context.Background(), bson.D{})
		if err != nil {
			return
		}
		for _, name := range names {
			test(db, name)
		}
	}
	pri()
}

func test(db *mongo.Database, coll string) {
	pipeline := mongo.Pipeline{
		bson.D{
			{"$collStats", bson.D{
				{"latencyStats", bson.D{
					//{"histogram", true},
				}},
			}},
		},
	}
	cursor, err := db.Collection(coll).Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// Iterate through the results
	var result []bson.M
	if err := cursor.All(context.TODO(), &result); err != nil {
		log.Fatal(err)
	}

	stats := result[0]["latencyStats"].(bson.M)
	readStats := stats["reads"].(bson.M)
	//klog.Infof("%s.%s => %+v \n", db.Name(), coll, readStats)

	reads = append(reads, read{
		latency: readStats["latency"].(int64),
		ops:     readStats["ops"].(int64),
		db:      db.Name(),
		coll:    coll,
	})
}

func pri() {
	sortByLatency := func(slice []read) {
		sort.Slice(slice, func(i, j int) bool {
			return slice[i].latency > slice[j].latency
		})
	}

	sortByLatency(reads)
	for _, r := range reads {
		klog.Infof("%s.%s >> %d %d \n", r.db, r.coll, r.latency, r.ops)
	}
}
