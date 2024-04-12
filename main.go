package main

import (
	"context"
	"github.com/ArnobKumarSaha/mongo/insert"
	"github.com/ArnobKumarSaha/mongo/mongoclient"
	"log"
)

func main() {
	client := mongoclient.ConnectFromPod()
	//client := mongoclient.ConnectLocal()
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	//object_count.Run(client)
	insert.Run(client)
	//latency.Run(client)
}
