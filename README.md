## Features
- latency
- object-count
- insert
- get-stats

## Run this
- from another pod
```bash
Edit the yamls/pod.yaml file to set the ENVs

make
```

- from Localhost
```bash
kubectl port-forward -n demo svc/source 27018:27017

export KUBECONFIG=/home/arnob/Downloads/configs/ui-demo-kubeconfig.yaml
export MONGODB_NAMESPACE=demo
export MONGODB_NAME=source
export INSERTION_SIZE_IN_GiB=5
```

## Observation
For 5Gi data insertion, The max usage touches 5 cpu, 5Gi memory & 15.5Gi storage. Takes 10m30s.
For 10Gi data insertion, The max usage touches 6.6 cpu, 6Gi memory & 24Gi storage. Takes 18m.
For 15Gi data insertion, The max usage touches 6.6 cpu, 6.4Gi memory & 35Gi storage. Takes 26m.
For 20Gi data insertion, The max usage touches 6.6 cpu, 6.4Gi memory & 45Gi storage. Takes 33m.

For 50Gi data insertion, The max usage touches 6.8 cpu, 8.6Gi memory & 108Gi storage. Takes 1h30m.
For 100Gi data insertion, The max usage touches  cpu, Gi memory & Gi storage. Takes .