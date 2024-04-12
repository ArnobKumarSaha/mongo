package latency

import (
	"context"
	"github.com/ArnobKumarSaha/mongo/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"k8s.io/klog/v2"
	"log"
	"sort"
)

func Run(client *mongo.Client) {
	collMap := database.ListCollectionsForAllDatabases(client)
	for db, collections := range collMap {
		for _, coll := range collections {
			CalculateLatency(client.Database(db), coll)
		}
	}
	PrintLatencyStats()
}

type read struct {
	latency, ops int64
	db, coll     string
}

var reads []read

func CalculateLatency(db *mongo.Database, coll string) {
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

func PrintLatencyStats() {
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
