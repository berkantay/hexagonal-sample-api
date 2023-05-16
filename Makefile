FILE_CONTENT := $(shell cat migrations/geofence_newyork.up.tile38)

build:
	go build cmd/api/main.go

service-up:
	echo "Run docker compose to init tile38, application, redis"

geofence-migrate-newyork-local: #this part could also be done using tile38-cli.
	curl -X POST 'localhost:9851' \
	-H 'Content-Type: text/plain' \
	--data-raw '$(FILE_CONTENT)'

.PHONY: build service-up geofence-migrate-newyork-local
