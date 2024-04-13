all: image
	kubectl apply -f yamls/rbac.yaml
	kubectl apply -f yamls/pod.yaml

run:
	kubectl delete -f yamls/pod.yaml || true
	$(MAKE) image
	kubectl apply -f yamls/pod.yaml

image:
	docker build -t arnobkumarsaha/mongo-util .
	docker push arnobkumarsaha/mongo-util
	#kind load docker-image arnobkumarsaha/mongo-util

clean:
	kubectl delete -f yamls/pod.yaml
	kubectl delete -f yamls/rbac.yaml