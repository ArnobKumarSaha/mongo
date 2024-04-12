package object_count

import (
	"context"
	"fmt"
	"github.com/ArnobKumarSaha/mongo/database"
	"github.com/ArnobKumarSaha/mongo/k8s"
	"github.com/ArnobKumarSaha/mongo/mongoclient"
	"go.mongodb.org/mongo-driver/mongo"
	"k8s.io/klog/v2"
	kubedb "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	"log"
	"time"
)

var (
	mg       *kubedb.MongoDB
	password string

	primaryPod   string
	secondaryPod string
	err          error
)

func Run(client *mongo.Client) {
	mg, err = mongoclient.GetMongoDB()
	if err != nil {
		klog.Fatal(err)
	}

	secret, err := mongoclient.GetSecret(mg.Spec.AuthSecret.Name, mg.Namespace)
	if err != nil {
		klog.Fatal(err)
	}
	password = string(secret.Data["password"])

	klog.Infof("MongoDB found : %v \n", mg.Name)
	hosts, err := database.GetPrimaryAndSecondaries(context.TODO(), client)
	if err != nil {
		_ = fmt.Errorf("error while getting primary and secondaries %v", err)
		return
	}
	primaryPod = hosts[0]
	secondaryPod = hosts[1]

	klog.Infof("Primary and Secondary found! %v %v \n", primaryPod, secondaryPod)

	tunnelPod, err := mongoclient.TunnelToDBPod(k8s.GetRESTConfig(), mg.Namespace, secondaryPod)
	if err != nil {
		return
	}

	klog.Infof("Tunnel created for pod %v at %v \n", secondaryPod, tunnelPod.Local)

	start := time.Now()
	klog.Infof("starts at %v \n", start)
	sc := mongoclient.ConnectToPod(tunnelPod, password)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	defer func() {
		if err := sc.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	_, err = CompareObjectsCount(client, sc)
	if err != nil {
		return
	}
	klog.Infof("Compares took %s", time.Since(start))
}
