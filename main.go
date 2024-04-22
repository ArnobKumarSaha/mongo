package main

import (
	"context"
	"github.com/ArnobKumarSaha/mongo/mongoclient"
	"github.com/ArnobKumarSaha/mongo/object_count"
	"github.com/ArnobKumarSaha/mongo/stats"
	"k8s.io/klog/v2"
	"log"
)

func main() {
	client := mongoclient.ConnectFromPod()
	//client := mongoclient.ConnectLocal()
	defer func() {
		klog.Infof("disconnecting in defer")
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	object_count.Run(client)
	//insert.Run(client)
	//latency.Run(client)
	stats.Run(client)
}
