package insert

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	wg           sync.WaitGroup
	stopPrinting chan struct{}
	start        time.Time

	dataSize      = 1024 * 1024 * 512
	batchSize     = 200000
	numGoroutines = 5
	databases     = []string{"aa", "bb", "cc"}
	collections   = []string{"one", "two", "three", "four"}
	checkpoints   []float64
)

func startInsertionGoRoutine(ctx context.Context, client *mongo.Client, rt int) {
	defer wg.Done()
	fmt.Printf("Starts goroutine %d.\n", rt)
	for {
		select {
		case <-stopPrinting:
			fmt.Println("Stopping goroutine.")
			return
		default:
			batch := make([]interface{}, batchSize)
			for b := 0; b < batchSize; b++ {
				document := map[string]interface{}{
					"key":   rand.Int() % 50,
					"value": rand.Float64(),
				}
				batch[b] = document
			}
			coll := getRandomCollection(client)
			_, err := coll.InsertMany(ctx, batch)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Inserted %d documents in %s/%s \n", batchSize, coll.Database().Name(), coll.Name())
		}
	}
}

func getRandomCollection(client *mongo.Client) *mongo.Collection {
	// Generate random index for databases and collections slices
	databaseIndex := rand.Intn(len(databases))
	collectionIndex := rand.Intn(len(collections))

	db := client.Database(databases[databaseIndex])
	collection := db.Collection(collections[collectionIndex])
	return collection
}

func Run(client *mongo.Client) {
	start = time.Now()
	stopPrinting = make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Hour)
	defer cancel()

	rand.Seed(time.Now().UnixNano())

	err := checkpointPreviousSizes(ctx, client)
	if err != nil {
		fmt.Printf("err on checkpointing = %v", err)
		return
	}
	go checkCondition(ctx, client)
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go startInsertionGoRoutine(ctx, client, i)
	}
	wg.Wait()
}
