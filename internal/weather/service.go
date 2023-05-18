package weather

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/berkantay/firefly-weather-condition-api/internal/index"
)

type CacheRepository interface {
	Set(ctx context.Context, cacheWeather *domain.Cache, ttl time.Duration) (*domain.Cache, error)
	Get(ctx context.Context, key string) (*domain.Cache, error)
	Exists(ctx context.Context, key string) bool
}

type GeospatialRepository interface {
	CityIntersect(ctx context.Context, coordinate *domain.Coordinate) bool
}

type WeatherClientRepository interface {
	FetchWeather(ctx context.Context, coordinate *domain.Coordinate) (*domain.Weather, error)
}

type WeatherService struct {
	cacheRepository         CacheRepository
	geospatialDatabase      GeospatialRepository
	weatherClientRepository WeatherClientRepository
}

func NewService(cache CacheRepository, geospatialDatabase GeospatialRepository, weatherClientRepository WeatherClientRepository) *WeatherService {
	return &WeatherService{
		cacheRepository:         cache,
		geospatialDatabase:      geospatialDatabase,
		weatherClientRepository: weatherClientRepository,
	}
}

func (ws *WeatherService) GetWeather(ctx context.Context, coordinate *domain.Coordinate) (*domain.Weather, error) {
	if !ws.geospatialDatabase.CityIntersect(ctx, coordinate) {
		return &domain.Weather{}, errors.New("the point is not in the market area")
	}

	key := index.CreatKey(coordinate.Latitude, coordinate.Longitude, 9) //TODO: resolution value could be adjusted 0-16 according to granularity demands

	if ws.cacheRepository.Exists(ctx, key) {
		var weather domain.Weather
		cache, err := ws.cacheRepository.Get(ctx, key)
		if err != nil {
			return nil, errors.New("cache get operation failed")
		}

		err = json.Unmarshal([]byte(cache.Value), &weather)
		if err != nil {
			return nil, err
		}
		return &weather, nil
	}

	weather, err := ws.weatherClientRepository.FetchWeather(ctx, coordinate)
	if err != nil {
		return nil, err
	}

	log.Println("Weather is", weather)

	weatherBytes, err := json.Marshal(weather)
	if err != nil {
		return nil, err
	}

	ws.cacheRepository.Set(ctx, &domain.Cache{
		Key:   key,
		Value: string(weatherBytes),
	}, 60*time.Second)

	return weather, nil
}
