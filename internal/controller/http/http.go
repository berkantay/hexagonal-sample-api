package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/gin-gonic/gin"
)

// WeatherService is an interface that defines methods for retrieving weather information.
type WeatherService interface {
	GetWeather(ctx context.Context, coordinate *domain.Coordinate) (*domain.Weather, error)
}

// WeatherHandler handles HTTP requests related to weather.
type WeatherHandler struct {
	WeatherService WeatherService
}

// NewWeatherHandler creates a new instance of WeatherHandler.
func NewWeatherHandler(engine *gin.Engine, weatherService WeatherService) {
	wh := &WeatherHandler{
		WeatherService: weatherService,
	}
	engine.GET("/weather", wh.GetWeather)
}

// GetWeather handles the GET /weather endpoint.
func (wh *WeatherHandler) GetWeather(c *gin.Context) {
	var coordinate domain.Coordinate

	err := c.ShouldBindQuery(&coordinate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})
		return
	}

	weather, err := wh.WeatherService.GetWeather(c.Request.Context(), &coordinate)
	if err != nil {
		switch {
		case err.Error() == "the point is not in the market area":
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, weather)
}
