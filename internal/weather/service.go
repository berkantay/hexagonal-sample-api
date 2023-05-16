package weather

import (
	"context"
	"log"
	"time"

	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
)

type CacheRepository interface {
	Set(ctx context.Context, cacheWeather domain.Cache, ttl time.Duration) (*domain.Cache, error)
	Get(ctx context.Context, key string) (*domain.Cache, error)
}

type GeospatialRepository interface {
	CityIntersectByCode(ctx context.Context, code string, coordinate *domain.Coordinate) bool
}

type WeatherService struct {
	cacheRepository    CacheRepository
	geospatialDatabase GeospatialRepository
}

func NewService(cache CacheRepository, geospatialDatabase GeospatialRepository) *WeatherService {
	return &WeatherService{
		cacheRepository:    cache,
		geospatialDatabase: geospatialDatabase,
	}
}

func (ws *WeatherService) GetWeather(ctx context.Context, coordinate *domain.Coordinate) *domain.Weather {
	if !ws.geospatialDatabase.CityIntersectByCode(ctx, "ny", coordinate) {
		log.Println("not in area")
	}
	return &domain.Weather{}
}
