package http

import (
	"context"
	"net/http"

	docs "github.com/berkantay/firefly-weather-condition-api/docs"
	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	docs.SwaggerInfo.BasePath = "/"

	engine.GET("/weather", wh.GetWeather)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// GetWeather handles the GET /weather endpoint.
// @BasePath /

// Firefly-weather-condition-api godoc
// @Summary Get weather condition for a coordinate
// @Description Get weather condition for a coordinate.
// @Accept json
// @Produce json
// @Param   latitude    query    string     true        "Latitude"
// @Param   longitude     query    string     true        "Longitude"
// @Success 200 {array} domain.Weather
// @Failure 400 {object} domain.ErrorResponse
// @Failure 422 {object} domain.ErrorResponse
// @Router /weather/ [get]
func (wh *WeatherHandler) GetWeather(c *gin.Context) {
	var coordinate domain.Coordinate

	err := c.ShouldBindQuery(&coordinate)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "invalid request"})
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
