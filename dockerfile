FROM golang:1.19-alpine AS builder

WORKDIR /app

RUN apk add build-base
RUN apk update && apk add bash ca-certificates git 

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o firefly-weather-condition-api cmd/api/main.go

FROM alpine:3.14

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/firefly-weather-condition-api .
# COPY --from=builder /app/resources resources


EXPOSE 8081

CMD ["/app/firefly-weather-condition-api"]