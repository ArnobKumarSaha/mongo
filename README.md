## Features
- latency
- object-count
- insert
- get-stats

## Run this
- from another pod
```bash
export KUBECONFIG=/home/arnob/Downloads/configs/ui-demo-kubeconfig.yaml
export MONGODB_NAMESPACE=demo
export MONGODB_NAME=source
```

- from Localhost
```bash
kubectl port-forward -n demo svc/source 27018:27017
```

