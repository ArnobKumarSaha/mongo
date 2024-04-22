package insert

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"k8s.io/klog/v2"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	wg           sync.WaitGroup
	stopPrinting chan struct{}
	start        time.Time

	dataSize      = 1024 * 1024 * 1024 // 1Gi
	batchSize     = 100000
	numGoroutines = 5
	databases     = []string{"aa", "bb", "cc"}
	collections   = []string{"one", "two"}
	checkpoints   []float64
)

func init() {
	iSz := os.Getenv("INSERTION_SIZE_IN_GiB")
	if iSz != "" {
		iSzInt, err := strconv.Atoi(iSz)
		if err != nil {
			log.Fatal(err)
		}
		dataSize = dataSize * iSzInt
	}

	bSz := os.Getenv("BATCH_SIZE")
	if bSz != "" {
		bSzInt, err := strconv.Atoi(bSz)
		if err != nil {
			log.Fatal(err)
		}
		batchSize = bSzInt
	}
	gr := os.Getenv("NUMBER_OF_GOROUTINE")
	if gr != "" {
		grInt, err := strconv.Atoi(gr)
		if err != nil {
			log.Fatal(err)
		}
		numGoroutines = grInt
	}
	klog.Infof("INSERTION_SIZE_IN_GiB=%d, BATCH_SIZE=%d, NUMBER_OF_GOROUTINE=%d\n", dataSize/(1024*1024*1024), batchSize, numGoroutines)
}

func startInsertionGoRoutine(ctx context.Context, client *mongo.Client, rt int) {
	defer wg.Done()
	klog.Infof("Starts goroutine %d.\n", rt)
	for {
		select {
		case <-stopPrinting:
			klog.Infoln("Stopping goroutine.")
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

			klog.Infof("Inserted %d documents in %s/%s \n", batchSize, coll.Database().Name(), coll.Name())
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
		klog.Infof("err on checkpointing = %v", err)
		return
	}
	go checkCondition(ctx, client)
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go startInsertionGoRoutine(ctx, client, i)
	}
	wg.Wait()
}
