build:
	go build cmd/api/main.go
service-up:
	echo "Run docker compose to init tile38, application, redis"

.PHONY:
	build service-up