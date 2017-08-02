PROJECT=mysagaproject
VERSION=v1

rootdir=$(realpath .)

run:
	dev_appserver.py app.yaml

deploy:
	gcloud app deploy --version $(VERSION) --project $(PROJECT)

prom:
	docker run -it -p 9090:9090 -v $(rootdir)/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
