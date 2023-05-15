# firefly-weather-condition-api

This weather condition api collects and returns the weather conditions at a given Market.

## To-Do

- Write a service that returns weather conditions at a given Market.
  - The service should take latitude and longitude as input.
  - The service should check if the given Point(lat, long) is in the given Market, e.g. New York.
  - Tile38 - Ultra Fast Geospatial Database & Geofencing Server should be used for this check.
  - You should populate new_york.geojson to Tile38 at setup.
  - If the given Point is not in the Market, return a proper response such as "The Point is not in the market area".
  - If the Point is in the market, first check the cache, use: Redis - in-memory data structure store, then call a free Weather API service, then cache and return the response.
  - You should use Uber's H3 resolution 8 index as the cache key and set TTL for the key such as 60 seconds. You can find the H3 hierarchical geospatial indexing system details at https://h3geo.org/.
  - You can use a free weather service like https://rapidapi.com/weatherapi/api/weatherapi-com or any other one that you prefer.
- The service should be written using The GO Programming Language.
- The service should be dockerized.
- Unit and Integration Tests should be written. Test coverage of more than 85% is expected.
- The Clean Architecture is preferred, but you can use any one that you prefer.
- You should share the project on the Github repo with AliCaner and Mecit.
- Demonstrating progressive development (e.g. not just a commit) and proper GitHub usage is a plus.

## Proposal Architecture

![Architecture Blueprint](resources/Firefly.png)
