package main

import (
	"context"
	"github.com/ArnobKumarSaha/mongo/insert"
	"github.com/ArnobKumarSaha/mongo/mongoclient"
	"log"
)

func main() {
	client := mongoclient.ConnectLocal()
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	insert.Run(client)
}
