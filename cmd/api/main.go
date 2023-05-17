package main

import (
	"context"
	"fmt"
	"os"

	"github.com/berkantay/firefly-weather-condition-api/config"
	"github.com/berkantay/firefly-weather-condition-api/internal/controller/http"
	"github.com/berkantay/firefly-weather-condition-api/internal/repository/api"
	"github.com/berkantay/firefly-weather-condition-api/internal/repository/db"
	"github.com/berkantay/firefly-weather-condition-api/internal/repository/tiles"
	"github.com/berkantay/firefly-weather-condition-api/internal/weather"
	"github.com/berkantay/firefly-weather-condition-api/pkg/log"
	"github.com/gin-gonic/gin"
)

const version = "v0.0.1"

func main() {
	fmt.Println("Version:", version)
	fmt.Println("Hello Firefly!")

	logger := log.NewLogger("api.log")
	defer logger.Close()

	config, err := config.NewConfig(context.Background(), "development", "resources")
	if err != nil {
		logger.Warn("could not read configuration")
		os.Exit(1)
	}

	geospatialClient, err := tiles.NewClient(config)
	if err != nil {
		logger.Warn("could not connect geospatial database")
		os.Exit(1)
	}

	cache, err := db.NewRedisStorage(config)
	if err != nil {
		logger.Warn("could not connect cache db")
		os.Exit(1)
	}

	weatherClient, err := api.NewWeatherClient(config)
	if err != nil {
		logger.Warn("could not create client")
		os.Exit(1)
	}

	weatherService := weather.NewService(cache, geospatialClient, weatherClient)

	webEngine := gin.Default()
	http.NewWeatherHandler(webEngine, weatherService)

	webEngine.Run("0.0.0.0:8081")
}
