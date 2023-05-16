package config

import (
	"context"

	"github.com/spf13/viper"
)

// Config provides methods for parsing configurations
type Config struct {
	Redis      Redis
	Tile38     Tile38
	WeatherApi WeatherApi
}

// Creates new configuration object with given configuration name and configuration path
func NewConfig(ctx context.Context, configName string, configPath string) (*Config, error) {
	v, err := readConfig(configName, configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}
	return config, nil
}

// Reads config from given name and path parameters.
func readConfig(configName, configPath string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	err := v.ReadInConfig()
	return v, err
}
