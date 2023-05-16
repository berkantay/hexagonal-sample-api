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

type weatherService struct {
	cacheRepository CacheRepository
}

func NewWeatherService(cache CacheRepository) *weatherService {
	return &weatherService{
		cacheRepository: cache,
	}
}
