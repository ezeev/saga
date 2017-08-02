ORG=ezeev
APP=lagunatpcom
VERSION=0.1
PACKAGE_PATH=github.com/ezeev/saga/lagunatp-service

# Start cloud proxy
# cloud_sql_proxy -instances=cloud-ninjaio:us-central1:cloud-ninja=tcp:3306 -credential_file=keys/cloud-ninjaio-113e76944ae4.json

clean:
	rm -f $(APP)

run:
	go install
	$(APP)

build: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(APP) $(PACKAGE_PATH)

docker: build
	@docker build -t $(ORG)/$(APP)\:$(VERSION) -f deploy/Dockerfile .
	@docker build -t $(ORG)/$(APP)\:latest -f deploy/Dockerfile .

docker-push: docker
	@docker push $(ORG)/$(APP)\:$(VERSION)
	@docker push $(ORG)/$(APP)\:latest

kube-delete:
	kubectl delete deployment $(APP)
	kubectl delete service $(APP)

kube-create: docker-push
	kubectl create -f deploy/deployment.yml
	kubectl expose deployment $(APP) --target-port=8081 --type=NodePort