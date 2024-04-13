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

kubectl cp demo/util:/app/all-stats /tmp/data # For copying the output of stats commands from pod
```

- from Localhost
```bash
kubectl port-forward -n demo svc/source 27018:27017

export MONGODB_USERNAME=root
export MONGODB_PASSWORD=12345
EXPORT MONGODB_LOCAL_PORT=27018

go run main.go
```

## Statistics
For 5Gi data insertion, The max usage touches 5 cpu, 5Gi memory & 15.5Gi storage. Takes 10m30s. <br>
For 10Gi data insertion, The max usage touches 6.6 cpu, 6Gi memory & 24Gi storage. Takes 18m. <br>
For 15Gi data insertion, The max usage touches 6.6 cpu, 6.4Gi memory & 35Gi storage. Takes 26m. <br>
For 20Gi data insertion, The max usage touches 6.6 cpu, 6.4Gi memory & 45Gi storage. Takes 33m. <br>
For 50Gi data insertion, The max usage touches 6.8 cpu, 8.6Gi memory & 108Gi storage. Takes 1h30m. <br>
For 100Gi data insertion, The max usage touches 6.8 cpu, 8.6Gi memory & 215Gi storage. Takes 2h54m. <br>

After 1h30m of the 100Gi data insertion, this is the resource usages:
```bash
arnob@msi ~> kubectl top pods -n demo
NAME       CPU(cores)   MEMORY(bytes)   
source-0   209m         6184Mi          
source-1   177m         5757Mi          
source-2   186m         5850Mi     
```