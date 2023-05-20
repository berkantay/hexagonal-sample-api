FILE_CONTENT := $(shell cat migrations/geofence_newyork.up.tile38)

build:
	go build cmd/api/main.go

service-up:
	docker-compose build --no-cache
	docker-compose up

unit-test:
	go test ./... -coverprofile cover.out
	go tool cover -html cover.out -o cover.html
	go tool cover -func cover.out | grep total:

integration-test:
	sudo chmod +x testing/integration.sh
	./testing/integration.sh

geofence-migrate-newyork-local: #this part could also be done using tile38-cli.
	curl -X POST 'localhost:9851' \
	-H 'Content-Type: text/plain' \
	--data-raw '$(FILE_CONTENT)'

.PHONY: build service-up geofence-migrate-newyork-local unit-test
