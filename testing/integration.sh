#!/bin/bash

#1 Run docker compose up
#2 Inject geojson
#3 Send requests with cURL

echo "Running integration tests 🔍"

echo "Checking and cleaning up old test environments 🧹"
docker rm $(sudo docker stop $(sudo docker ps -a | grep "weather-api" | cut -d " " -f 1))
docker rm $(sudo docker stop $(sudo docker ps -a | grep "redis-cache" | cut -d " " -f 1))
docker rm $(sudo docker stop $(sudo docker ps -a | grep "tile38" | cut -d " " -f 1))



echo "Igniting test environment 🔥"
docker-compose build --no-cache
docker-compose up -d
echo "Injecting geojson into tile38 🗺️"
make geofence-migrate-newyork-local

status_code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8081/weather?latitude=40.731328&longitude=-74.067534")
# Output the status code
echo "Status code of correct query parameters: $status_code"
if [[ "$status_code" -ne 200 ]] ; then
  echo "integration test failed ❌ - on correct query parameters"
else
  echo "correct query parameters ✅"
fi


# Output the status code
status_code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8081/weather?latitude=40.731328&longitude=4.067534")
echo "Status code of non overlapping coordinate: $status_code"
if [[ "$status_code" -ne 422 ]] ; then
  echo "integration test failed ❌ - on non overlapping coordinate"
else
  echo "non overlapping coordinate ✅"
fi

# Output the status code
status_code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8081/weather")
echo "Status code of missing query parameters: $status_code"
if [[ "$status_code" -ne 400 ]] ; then
  echo "integration test failed ❌ - on missing query parameters"
else
  echo "missing query parameters ✅"
  exit 0
fi