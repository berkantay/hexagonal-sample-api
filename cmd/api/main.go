package main

import (
	"context"
	"fmt"
	"os"

	"github.com/berkantay/firefly-weather-condition-api/config"
	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/berkantay/firefly-weather-condition-api/internal/repository/tiles"
	"github.com/berkantay/firefly-weather-condition-api/pkg/log"
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

	fmt.Println(config)

	tileClient, err := tiles.NewClient(config)
	if err != nil {
		logger.Warn("could not connect geospatial database")
		os.Exit(1)
	}

	tileClient.CityIntersect(context.TODO(), domain.Coordinate{})
}
