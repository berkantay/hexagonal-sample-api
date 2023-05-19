package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/berkantay/firefly-weather-condition-api/config"
	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
)

// WeatherClient is responsible for fetching weather data from an API.
type WeatherClient struct {
	address       string
	XRapidAPIKey  string
	XRapidAPIHost string
}

// NewWeatherClient creates a new instance of WeatherClient.
func NewWeatherClient(config *config.Config) (*WeatherClient, error) {
	return &WeatherClient{
		address:       config.WeatherApi.Address,
		XRapidAPIKey:  config.WeatherApi.APIKey,
		XRapidAPIHost: config.WeatherApi.APIHost,
	}, nil
}

// FetchWeather retrieves weather data for the given coordinate.
func (wc *WeatherClient) FetchWeather(ctx context.Context, coordinate *domain.Coordinate) (*domain.Weather, error) {
	var weather domain.Weather

	queryString := buildCoordinateQuery(coordinate)

	params := url.Values{}
	params.Add("q", queryString)

	url := wc.address + params.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", wc.XRapidAPIKey)
	req.Header.Add("X-RapidAPI-Host", wc.XRapidAPIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Request error is:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	log.Println("Response is", string(body))
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, err
	}
	return &weather, nil
}

// buildCoordinateQuery creates a formatted query string for the given coordinate.
func buildCoordinateQuery(coordinate *domain.Coordinate) string {
	return fmt.Sprintf("%.2f,%.2f", coordinate.Latitude, coordinate.Longitude)
}
