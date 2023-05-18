package weather

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/berkantay/firefly-weather-condition-api/internal/index"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

type mockCacheRepository struct {
	setError  error
	getError  error
	mockCache *domain.Cache
	existBool bool
}

func (mcr *mockCacheRepository) Set(ctx context.Context, cacheWeather *domain.Cache, ttl time.Duration) (*domain.Cache, error) {
	return mcr.mockCache, mcr.setError
}

func (mcr *mockCacheRepository) Get(ctx context.Context, key string) (*domain.Cache, error) {
	return mcr.mockCache, mcr.getError
}

func (mcr *mockCacheRepository) Exists(ctx context.Context, key string) bool {
	return mcr.existBool
}

type mockGeospatialRepository struct {
	isIntersect bool
}

func (mgr *mockGeospatialRepository) CityIntersect(ctx context.Context, coordinate *domain.Coordinate) bool {
	return mgr.isIntersect
}

type mockWeatherClient struct {
	mockWeather *domain.Weather
	mockError   error
}

func (mwc *mockWeatherClient) FetchWeather(ctx context.Context, coordinate *domain.Coordinate) (*domain.Weather, error) {
	fmt.Println("Coordinates are", coordinate)
	return mwc.mockWeather, mwc.mockError
}

func TestCityDoesNotIntersect(t *testing.T) {
	t.Run("Given coordinates", func(t *testing.T) {
		mockCoordinate := domain.Coordinate{
			Latitude:  generateRandomLatitude(),
			Longitude: generateRandomLongitude(),
		}
		t.Run("When coordinates ARE NOT in newyork", func(t *testing.T) {
			mockGeo := mockGeospatialRepository{
				isIntersect: false,
			}
			t.Run("Then it should return error message", func(t *testing.T) {
				weatherService := NewService(&mockCacheRepository{}, &mockGeo, &mockWeatherClient{})
				result, err := weatherService.GetWeather(context.Background(), &mockCoordinate)
				assert.Equal(t, result, &domain.Weather{})
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "the point is not in the market area")

			})
		})
	})
}

func TestCityDoesIntersect(t *testing.T) {
	t.Run("Given coordinates", func(t *testing.T) {
		mockLat := generateRandomLatitude()
		mockLon := generateRandomLongitude()
		mockCoordinate := domain.Coordinate{
			Latitude:  mockLat,
			Longitude: mockLon,
		}
		t.Run("When coordinates does not in newyork", func(t *testing.T) {
			mockGeo := mockGeospatialRepository{
				isIntersect: true,
			}
			t.Run("Then it should return error message", func(t *testing.T) {
				weatherService := NewService(&mockCacheRepository{}, &mockGeo, &mockWeatherClient{})
				result, err := weatherService.GetWeather(context.Background(), &mockCoordinate)
				assert.NotEqual(t, result, &domain.Weather{})
				assert.Nil(t, err)

			})
		})
	})
}

func TestCityDoesIntersectExistInCacheButGetFail(t *testing.T) {
	t.Run("Given coordinates", func(t *testing.T) {
		mockLat := generateRandomLatitude()
		mockLon := generateRandomLongitude()
		mockCoordinate := domain.Coordinate{
			Latitude:  mockLat,
			Longitude: mockLon,
		}
		t.Run("When coordinates ARE in newyork, cached but get failed", func(t *testing.T) {
			mockGeo := &mockGeospatialRepository{
				isIntersect: true,
			}
			mockCache := &mockCacheRepository{
				getError:  errors.New("cache get operation failed"),
				setError:  nil,
				existBool: true,
				mockCache: &domain.Cache{},
			}
			t.Run("Then it should return error message", func(t *testing.T) {
				weatherService := NewService(mockCache, mockGeo, &mockWeatherClient{})
				result, err := weatherService.GetWeather(context.Background(), &mockCoordinate)
				assert.Nil(t, result)
				assert.Equal(t, err.Error(), mockCache.getError.Error())
			})
		})
	})
}

func TestCityDoesIntersectExistInCacheGetSuccess(t *testing.T) {
	t.Run("Given coordinates", func(t *testing.T) {
		mockLat := generateRandomLatitude()
		mockLon := generateRandomLongitude()
		mockCoordinate := domain.Coordinate{
			Latitude:  mockLat,
			Longitude: mockLon,
		}
		t.Run("When coordinates ARE in newyork, cache get success", func(t *testing.T) {

			tempH3Key := index.CreatKey(mockLat, mockLon, 9)

			mockGeo := &mockGeospatialRepository{
				isIntersect: true,
			}
			mockWeatherClientTest := &mockWeatherClient{
				mockWeather: &domain.Weather{
					Location: domain.Location{
						Name: "test",
					},
				},
				mockError: nil,
			}

			mockValue, err := json.Marshal(&domain.Weather{
				Location: domain.Location{
					Name: "test",
				},
			})
			if err != nil {
				log.Fatal(err)
			}

			mockCache := &mockCacheRepository{
				getError:  nil,
				setError:  nil,
				existBool: true,
				mockCache: &domain.Cache{
					Key:   tempH3Key,
					Value: string(mockValue),
				},
			}
			t.Run("Then it should return weather", func(t *testing.T) {
				weatherService := NewService(mockCache, mockGeo, mockWeatherClientTest)
				result, err := weatherService.GetWeather(context.Background(), &mockCoordinate)
				assert.Nil(t, err)
				assert.Equal(t, mockCache.mockCache.Key, tempH3Key)
				fmt.Println(result.Location)
				assert.Equal(t, mockWeatherClientTest.mockWeather.Location.Name, result.Location.Name)
			})
		})
	})
}

func generateRandomLatitude() float64 {
	// Latitude ranges from -90 to 90
	min := -90.0
	max := 90.0

	// Generate a random float64 value between min and max
	return rand.Float64()*(max-min) + min
}

func generateRandomLongitude() float64 {
	// Longitude ranges from -180 to 180
	min := -180.0
	max := 180.0

	// Generate a random float64 value between min and max
	return rand.Float64()*(max-min) + min
}
