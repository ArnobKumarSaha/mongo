package insert

import (
	"context"
	"fmt"
	"github.com/ArnobKumarSaha/mongo/database"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"time"
)

func checkCondition(ctx context.Context, client *mongo.Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sz, _ := calcTotalStorageSize(ctx, client)
			if int(sz) >= dataSize {
				fmt.Printf("Stopping all goroutines. Time elapsed: %v\n", time.Since(start))
				close(stopPrinting)
				return
			}
			convertToReadableUnit(sz)
		}
	}
}

func checkpointPreviousSizes(ctx context.Context, client *mongo.Client) error {
	checkpoints = make([]float64, len(databases))
	for i := 0; i < len(databases); i++ {
		databaseName := databases[i]
		dbStats, err := database.DBStats(ctx, client, databaseName)
		if err != nil {
			return err
		}
		storageSz, ok := dbStats["storageSize"]
		if !ok {
			return fmt.Errorf("type assertion error: can't get dataSize info")
		}

		switch reflect.TypeOf(storageSz).Name() {
		case "float64":
			checkpoints[i] = storageSz.(float64)
			break
		case "int32":
			checkpoints[i] = float64(storageSz.(int32))
			break
		default:
			_ = fmt.Errorf("%v type is invalid", reflect.TypeOf(storageSz).Name())
		}

	}
	return nil
}

func calcTotalStorageSize(ctx context.Context, client *mongo.Client) (float64, error) {
	totalSize := float64(0)
	for i := 0; i < len(databases); i++ {
		databaseName := databases[i]
		dbStats, err := database.DBStats(ctx, client, databaseName)
		if err != nil {
			return 0, err
		}

		storageSz, ok := dbStats["storageSize"]
		if !ok {
			return 0, fmt.Errorf("type assertion error: can't get dataSize info")
		}

		totalSize += storageSz.(float64) - checkpoints[i]
	}
	return totalSize, nil
}

func convertToReadableUnit(sz float64) {
	divCount := 0
	for {
		if sz >= 1024 {
			sz = sz / 1024
			divCount += 1
		} else {
			break
		}
	}
	unit := "Byte"
	switch divCount {
	case 1:
		unit = "KiB"
		break
	case 2:
		unit = "MiB"
		break
	case 3:
		unit = "GiB"
		break
	case 4:
		unit = "TiB"
		break
	default:
		fmt.Printf("divCount %v is way too big \n", divCount)
	}
	fmt.Printf("%v %s inserted. Current time %v \n", sz, unit, time.Now())
}
