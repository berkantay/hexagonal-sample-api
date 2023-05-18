package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockWeatherService struct {
	mockWeather *domain.Weather
	mockError   error
}

func (mws *mockWeatherService) GetWeather(ctx context.Context, coordinate *domain.Coordinate) (*domain.Weather, error) {
	return mws.mockWeather, mws.mockError
}

func TestWeatherHandlerGetWeatherCorrect(t *testing.T) {
	t.Run("Given the http server is running", func(t *testing.T) {
		router := gin.Default()
		weatherService := &mockWeatherService{
			mockWeather: &domain.Weather{
				Location: domain.Location{
					Name: "los-angeles",
				},
				Current: domain.Current{
					TempC: 35.6,
				},
			},
		}
		NewWeatherHandler(router, weatherService)
		t.Run("When GET request is sent with CORRECT query parameters", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/weather?latitude=40.7128&longitude=-74.0060", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			t.Run("Then response must be correct ", func(t *testing.T) {
				assert.Equal(t, http.StatusOK, rec.Code)
				fmt.Println("Result is", rec.Body)
				var resultWeather domain.Weather
				json.Unmarshal(rec.Body.Bytes(), &resultWeather)
				assert.Equal(t, "los-angeles", resultWeather.Location.Name)
			})
		})
	})
}

func TestWeatherHandlerGetWeatherWrong(t *testing.T) {
	t.Run("Given the http server is running", func(t *testing.T) {
		router := gin.Default()
		weatherService := &mockWeatherService{
			mockWeather: nil,
			mockError:   errors.New("test"),
		}
		NewWeatherHandler(router, weatherService)
		t.Run("When GET request is sent with WRONG query parameters", func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/weather?latitude=&longitude=-74.0060", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			t.Run("Then response must return error ", func(t *testing.T) {
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				fmt.Println("Result is", rec.Body)
			})
		})
	})
}
