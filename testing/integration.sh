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

response=$(curl --location --request GET "http://localhost:8081/weather?latitude=40.331328&longitude=-74.077534" -s -D -)
status_code=$(echo "$response" | grep -i '^HTTP' | awk '{print $2}')
response_body=$(echo "$response" | sed -e '1,/^\r$/d')

if [[ "$status_code" -ne 200 ]] ; then
  echo "integration test failed ❌ "
else
  echo "correct query parameters ✅"
fi

#Response body match could not performed because external api and wheather condition is changed by time. Therefore expected body cannot be statically written

response=$(curl --location --request GET "http://localhost:8081/weather?latitude=40.731328&longitude=4.067534" -s -D -)
status_code=$(echo "$response" | grep -i '^HTTP' | awk '{print $2}')
response_body=$(echo "$response" | sed -e '1,/^\r$/d')
expected_response='{"error":"the point is not in the market area"}'

if [[ "$status_code" -ne 422 ]] ; then
  echo "integration test failed - on non overlapping coordinate❌ "
else
  echo "non overlapping coordinate ✅"
fi

if [[ "$response_body" == "$expected_response" ]]; then
  echo "Response body matches the expected JSON.-- non overlapping coordinate ✅ "
else
  echo "Response body DOES NOT matches the expected JSON. --  non overlapping coordinate" ❌
fi

response=$(curl --location --request GET "http://localhost:8081/weather" -s -D -)
status_code=$(echo "$response" | grep -i '^HTTP' | awk '{print $2}')
response_body=$(echo "$response" | sed -e '1,/^\r$/d')
expected_response='{"error":{}}'
if [[ "$status_code" -ne 400 ]] ; then
  echo "integration test failed -- on missing query parameters ❌ "
else
  echo "missing query parameters ✅"
fi

if [[ "$response_body" == "$expected_response" ]]; then
  echo "Response body matches the expected JSON. -- on missing query parameters ✅ "
else
  echo "Response body DOES NOT matches the expected JSON. -- on missing query parameters ❌ "
fi