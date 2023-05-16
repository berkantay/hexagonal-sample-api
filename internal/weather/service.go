package weather

import (
	"context"
	"time"

	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
)

type CacheRepository interface {
	Set(ctx context.Context, cacheWeather domain.Cache, ttl time.Duration) (*domain.Cache, error)
	Get(ctx context.Context, key string) (*domain.Cache, error)
}

type GeospatialRepository interface {
	CheckArea(coordinate domain.Coordinate)
}

type WeatherService struct {
	cacheRepository CacheRepository
}

func NewService(cache CacheRepository) *WeatherService {
	return &WeatherService{
		cacheRepository: cache,
	}
}

func (ws *WeatherService) GetWeather(ctx context.Context, coordinate *domain.Coordinate) *domain.Weather {
	return &domain.Weather{}
}
