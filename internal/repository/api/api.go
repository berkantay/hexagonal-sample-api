package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/berkantay/firefly-weather-condition-api/config"
	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
)

type WeatherClient struct {
	address       string
	XRapidAPIKey  string
	XRapidAPIHost string
}

func NewWeatherClient(config *config.Config) (*WeatherClient, error) {
	return &WeatherClient{
		address:       config.WeatherApi.Address,
		XRapidAPIKey:  config.WeatherApi.APIKey,
		XRapidAPIHost: config.WeatherApi.APIHost,
	}, nil
}

func (wc *WeatherClient) Fetch(ctx context.Context) (*domain.Weather, error) {
	var weather domain.Weather
	req, _ := http.NewRequest("GET", wc.address, nil)

	req.Header.Add("X-RapidAPI-Key", wc.XRapidAPIKey)
	req.Header.Add("X-RapidAPI-Host", wc.XRapidAPIHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, err
	}

	return &weather, nil
}
