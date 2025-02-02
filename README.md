# Hexagonal architecture example

This weather condition api collects and returns the weather conditions at a given Market.

## Testing

```
make unit-test
```

```
make integration-test
```

Then total coverage is %86.4. Also cover.html is generated to check the covered and uncovered lines.

## Usage

```
docker-compose build --no-cache
docker-compose up
```

`docker-compose up -d` if daemon mode desired.
Then application will be running at **localhost:8081**

# Important

To inject polygon into _Tile38_

```
make geofence-migrate-newyork-local
```

## API Documentation

API Documentation implemented. Visit `http://localhost:8081/swagger/index.html`

### Get temperature on given coordinate

#### Request

`GET /weather/`

`curl "http://localhost:8081/weather?latitude=40.731328&longitude=-74.067534"`

### Response - Success

```
{
  "location": {
    "name": "Jersey City",
    "region": "New Jersey",
    "country": "United States of America",
    "lat": 40.73,
    "lon": -74.07,
    "tz_id": "America/New_York",
    "localtime_epoch": 1684305567,
    "localtime": "2023-05-17 2:39"
  },
  "current": {
    "last_updated_epoch": 1684305000,
    "last_updated": "2023-05-17 02:30",
    "temp_c": 20.6,
    "temp_f": 69.1,
    "is_day": 0,
    "condition": {
      "text": "Clear",
      "icon": "//cdn.weatherapi.com/weather/64x64/night/113.png",
      "code": 1000
    },
    "wind_mph": 2.2,
    "wind_kph": 3.6,
    "wind_degree": 327,
    "wind_dir": "NNW",
    "pressure_mb": 1001,
    "pressure_in": 29.56,
    "precip_mm": 0,
    "precip_in": 0,
    "humidity": 39,
    "cloud": 0,
    "feelslike_c": 20.6,
    "feelslike_f": 69.1,
    "vis_km": 16,
    "vis_miles": 9,
    "uv": 1,
    "gust_mph": 17,
    "gust_kph": 27.4
  }
}
```

### Response - Point not in the market area

```
{
    "error": "the point is not in the market area"
}
```
