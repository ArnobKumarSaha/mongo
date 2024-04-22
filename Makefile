all: image
	kubectl apply -f yamls/rbac.yaml
	kubectl apply -f yamls/job.yaml

run:
	kubectl delete -f yamls/job.yaml || true
	$(MAKE) image
	kubectl apply -f yamls/job.yaml

image:
	docker build -t arnobkumarsaha/mongo-util .
	docker push arnobkumarsaha/mongo-util
	#kind load docker-image arnobkumarsaha/mongo-util

clean:
	kubectl delete -f yamls/job.yaml
	kubectl delete -f yamls/rbac.yaml