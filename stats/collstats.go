package stats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ArnobKumarSaha/mongo/database"
	"github.com/ArnobKumarSaha/mongo/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"path/filepath"
	"strings"
)

const dir = "collected-stats"

func Run(client *mongo.Client) {
	utils.MakeDir(dir)
	collMap := database.ListCollectionsForAllDatabases(client)
	for db, collections := range collMap {
		utils.MakeDir(filepath.Join(dir, db))
		dbRef := client.Database(db)
		databaseStats(dbRef)
		for _, coll := range collections {
			collectionStats(dbRef, coll)
		}
	}
}

func databaseStats(db *mongo.Database) {
	cmd := bson.D{{"dbStats", 1}, {"scale", 1048576}}
	var result bson.M
	err := db.RunCommand(context.TODO(), cmd).Decode(&result)
	if err != nil {

	}
	indentedData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	utils.WriteFile(filepath.Join(dir, db.Name()), "_", indentedData)
}

func collectionStats(db *mongo.Database, coll string) {
	cmd := bson.D{{"collStats", coll}}
	var result bson.M
	err := db.RunCommand(context.TODO(), cmd).Decode(&result)
	if err != nil {
		if strings.Contains(err.Error(), "is a view, not a collection") {
			fmt.Println(err.Error())
			return
		} else {
			log.Fatal(err)
		}
	}

	//var b []byte
	//buf := bytes.NewBuffer(b)
	//enc := json.NewEncoder(buf)
	//err = enc.Encode(&result)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//writeFile(fmt.Sprintf("%s.%s", db.Name(), coll), buf.Bytes())

	indentedData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	utils.WriteFile(filepath.Join(dir, db.Name()), coll, indentedData)
	// fmt.Sprintf("%s.%s", db.Name(), coll)
}
