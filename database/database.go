package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ListDatabases(client *mongo.Client) []string {
	dbs, err := client.ListDatabaseNames(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	return dbs
}

func ListCollectionsForSpecificDatabase(client *mongo.Client, database string) []string {
	db := client.Database(database)
	names, err := db.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}
	return names
}

func ListCollectionsForAllDatabases(client *mongo.Client) map[string][]string {
	dbs := ListDatabases(client)
	mp := make(map[string][]string)
	for _, d := range dbs {
		db := client.Database(d)
		names, err := db.ListCollectionNames(context.Background(), bson.D{})
		if err != nil {
			panic(err)
		}
		mp[d] = names
	}
	return mp
}

func DBStats(ctx context.Context, client *mongo.Client, db string) (map[string]interface{}, error) {
	dbStats := make(map[string]interface{})
	err := client.Database(db).RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}}).Decode(&dbStats)
	return dbStats, err
}
