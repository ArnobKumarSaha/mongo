all:
	#export CGO_ENABLED=0
	#go build
	docker build -t arnobkumarsaha/mongo-util .
	docker push arnobkumarsaha/mongo-util
	#kind load docker-image arnobkumarsaha/mongo-util
	kubectl apply -f yamls/rbac.yaml
	kubectl apply -f yamls/pod.yaml

run:
	kubectl delete -f yamls/pod.yaml || true
	docker build -t arnobkumarsaha/mongo-util .
	docker push arnobkumarsaha/mongo-util
	kubectl apply -f yamls/pod.yaml

clean:
	kubectl delete -f yamls/pod.yaml
	kubectl delete -f yamls/rbac.yaml