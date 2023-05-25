package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/berkantay/firefly-weather-condition-api/config"
	"github.com/berkantay/firefly-weather-condition-api/internal/controller/http"
	"github.com/berkantay/firefly-weather-condition-api/internal/repository/api"
	"github.com/berkantay/firefly-weather-condition-api/internal/repository/db"
	"github.com/berkantay/firefly-weather-condition-api/internal/repository/tiles"
	"github.com/berkantay/firefly-weather-condition-api/internal/weather"
	"github.com/berkantay/firefly-weather-condition-api/pkg/log"
	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
)

// Version indicates the current version of the application.
var Version = "development"

func main() {
	banner := fmt.Sprintf("Firefly-%s", Version)
	bannerFigure := figure.NewColorFigure(banner, "doom", "white", true)
	bannerFigure.Print()

	logger := log.NewLogger("api.log")
	defer logger.Close()

	config, err := config.NewConfig(context.Background(), "development", "resources")
	if err != nil {
		logger.Warn("could not read configuration, checking environment variables", map[string]interface{}{
			"error": err,
		})
		configureFromEnvironment(config)
		fmt.Println("Could not read configuration, checking environment variables...")
	}

	geospatialClient, err := tiles.NewClient(config)
	if err != nil {
		logger.Warn("could not connect geospatial database", map[string]interface{}{
			"error": err,
		})
		fmt.Println("could not connect geospatial database")

	}

	cache, err := db.NewRedisStorage(config)
	if err != nil {
		logger.Warn("could not connect cache db", map[string]interface{}{
			"error": err,
		})
		fmt.Println("could not connect cache db")
	}

	weatherClient, err := api.NewWeatherClient(config)
	if err != nil {
		logger.Warn("could not create client", map[string]interface{}{
			"error": err,
		})
		fmt.Println("could not create client")

	}

	weatherService := weather.NewService(logger, cache, geospatialClient, weatherClient)

	webEngine := gin.Default()
	http.NewWeatherHandler(webEngine, weatherService)

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-gracefulShutdown // Wait for the signal
		logger.Info("Shutting down...")
		cache.Close(context.TODO())
		geospatialClient.Close(context.TODO())
		os.Exit(0)
	}()

	// Start the web server
	err = webEngine.Run("0.0.0.0:8081")
	if err != nil {
		logger.Warn("failed to start server")
	}
}

func configureFromEnvironment(conf *config.Config) {
	conf.Redis.Host = os.Getenv("REDIS_HOST")
	conf.Redis.Port = os.Getenv("REDIS_PORT")
	conf.Redis.Password = os.Getenv("REDIS_PASSWORD")

	conf.Tile38.Host = os.Getenv("TILE38_HOST")
	conf.Tile38.Port = os.Getenv("TILE38_PORT")

	conf.WeatherApi.Address = os.Getenv("WEATHER_API_ADDRESS")
	conf.WeatherApi.APIHost = os.Getenv("WEATHER_API_HOST")
	conf.WeatherApi.APIKey = os.Getenv("WEATHER_API_KEY")
}
