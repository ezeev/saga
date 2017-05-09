PROJECT=mysagaproject
VERSION=v1

run:
	dev_appserver.py app.yaml

deploy:
	gcloud app deploy --version $(VERSION) --project $(PROJECT)
