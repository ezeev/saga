APP=viroonga-wp

# Start cloud proxy
# cloud_sql_proxy -instances=cloud-ninjaio:us-central1:cloud-ninja=tcp:3306 -credential_file=keys/cloud-ninjaio-113e76944ae4.json

kube-delete:
	kubectl delete deployment $(APP)
	kubectl delete service $(APP)

kube-create:
	kubectl create -f deploy.yml
	kubectl expose deployment $(APP) --target-port=80 --type=NodePort
	kubectl expose deployment $(APP) --target-port=80 --type=LoadBalancer --name=temp-lb
