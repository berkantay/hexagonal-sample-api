package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/gin-gonic/gin"
)

type WeatherService interface {
	GetWeather(ctx context.Context, coordinate *domain.Coordinate) (*domain.Weather, error)
}

type WeatherHandler struct {
	WeatherService WeatherService
}

func NewWeatherHandler(engine *gin.Engine, weatherService WeatherService) {
	wh := &WeatherHandler{
		WeatherService: weatherService,
	}
	engine.GET("/weather", wh.GetWeather)
}

func (wh *WeatherHandler) GetWeather(c *gin.Context) {
	var coordinate domain.Coordinate

	err := c.ShouldBindQuery(&coordinate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid request")})
		return
	}

	weather, err := wh.WeatherService.GetWeather(c.Request.Context(), &coordinate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, weather)

}
