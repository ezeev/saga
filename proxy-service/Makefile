ORG=ezeev
APP=nginx
VERSION=0.1
PACKAGE_PATH=github.com/ezeev/saga/proxy

docker:
	@docker build -t $(ORG)/$(APP)\:$(VERSION) -f deploy/Dockerfile .
	@docker build -t $(ORG)/$(APP)\:latest -f deploy/Dockerfile .

docker-run: docker
	@docker run -it -p 80:80 $(ORG)/$(APP):$(VERSION)

docker-push: docker
	@docker push $(ORG)/$(APP)\:$(VERSION)
	@docker push $(ORG)/$(APP)\:latest

kube-delete:
	kubectl delete deployment $(APP)
	kubectl delete service $(APP)

kube-update-image: docker-push
	kubectl set image deployment/$(APP) $(APP)=$(ORG)/$(APP):$(VERSION)

kube-create: docker-push
	kubectl create -f deploy/kube-deploy.yml
	kubectl expose deployment $(APP) --target-port=80 --type=NodePort