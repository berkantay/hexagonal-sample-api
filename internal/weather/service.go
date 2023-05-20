package weather

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/berkantay/firefly-weather-condition-api/internal/index"
	"github.com/berkantay/firefly-weather-condition-api/pkg/log"
)

// CacheRepository provides methods for caching weather data.
type CacheRepository interface {
	// Set stores the given cache weather data in the cache.
	// It takes a context, cacheWeather object, and time-to-live (TTL) duration as parameters.
	// It returns the stored cache weather data or an error if the operation fails.
	Set(ctx context.Context, cacheWeather *domain.Cache, ttl time.Duration) (*domain.Cache, error)
	// Get retrieves the cache weather data associated with the specified key from the cache.
	// It takes a context and key as parameters.
	// It returns the cache weather data or an error if the data retrieval fails.
	Get(ctx context.Context, key string) (*domain.Cache, error)
	// Exists checks if the cache weather data with the specified key exists in the cache.
	// It takes a context and key as parameters.
	// It returns true if the data exists, false otherwise.
	Exists(ctx context.Context, key string) bool
}

// GeospatialRepository provides methods for geospatial operations.
type GeospatialRepository interface {
	// CityIntersect checks if the given coordinate intersects with a city.
	// It takes a context and coordinate as parameters.
	// It returns true if the coordinate intersects with a city, false otherwise.
	CityIntersect(ctx context.Context, coordinate *domain.Coordinate) bool
}

// WeatherClientRepository provides methods for fetching weather data.
type WeatherClientRepository interface {
	// FetchWeather retrieves the weather data for the given coordinate.
	// It takes a context and coordinate as parameters.
	// It returns the weather data or an error if the data retrieval fails.
	FetchWeather(ctx context.Context, coordinate *domain.Coordinate) (*domain.Weather, error)
}

// WeatherService handles weather-related operations.
type WeatherService struct {
	cacheRepository         CacheRepository
	geospatialDatabase      GeospatialRepository
	weatherClientRepository WeatherClientRepository
	logger                  *log.Logger
}

// NewService creates a new instance of WeatherService.
func NewService(logger *log.Logger, cache CacheRepository, geospatialDatabase GeospatialRepository, weatherClientRepository WeatherClientRepository) *WeatherService {
	return &WeatherService{
		logger:                  logger,
		cacheRepository:         cache,
		geospatialDatabase:      geospatialDatabase,
		weatherClientRepository: weatherClientRepository,
	}
}

// GetWeather retrieves the weather data for the given coordinate.
func (ws *WeatherService) GetWeather(ctx context.Context, coordinate *domain.Coordinate) (*domain.Weather, error) {
	if !ws.geospatialDatabase.CityIntersect(ctx, coordinate) {
		ws.logger.Info("the point is not in the market area")
		return &domain.Weather{}, errors.New("the point is not in the market area")
	}

	key := index.CreateKey(coordinate.Latitude, coordinate.Longitude, 9) //TODO: resolution value could be adjusted 0-16 according to granularity demands

	if ws.cacheRepository.Exists(ctx, key) {
		var weather domain.Weather
		cache, err := ws.cacheRepository.Get(ctx, key)
		if err != nil {
			ws.logger.Warn("cache get operation failed")
			return nil, errors.New("cache get operation failed")
		}

		err = json.Unmarshal([]byte(cache.Value), &weather)
		if err != nil {
			ws.logger.Fatal(err.Error())
			return nil, err
		}
		return &weather, nil
	}

	weather, err := ws.weatherClientRepository.FetchWeather(ctx, coordinate)
	if err != nil {
		ws.logger.Fatal(err.Error())
		return nil, err
	}

	weatherBytes, err := json.Marshal(weather)
	if err != nil {
		ws.logger.Fatal(err.Error())
		return nil, err
	}

	ws.cacheRepository.Set(ctx, &domain.Cache{
		Key:   key,
		Value: string(weatherBytes),
	}, 60*time.Second)

	return weather, nil
}
