package tiles

import (
	"context"
	"fmt"

	"github.com/berkantay/firefly-weather-condition-api/config"
	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/xjem/t38c"
)

// Client is responsible for interacting with a Tile38 server.
type Client struct {
	client *t38c.Client
}

// NewClient creates a new instance of Client.
func NewClient(config *config.Config) (*Client, error) {
	tileUrl := fmt.Sprintf("%s:%s", config.Tile38.Host, config.Tile38.Port)
	client, err := t38c.New(t38c.Config{
		Address: tileUrl,
		Debug:   false, // print queries to stdout
	})
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		client: client,
	}, nil
}

// CityIntersect checks if a given coordinate intersects with a city in Tile38.
func (c *Client) CityIntersect(ctx context.Context, coordinate *domain.Coordinate) bool {
	data, err := c.client.Search.Intersects("cities").Circle(coordinate.Latitude, coordinate.Longitude, 0).Do(ctx)
	if err != nil {
		return false
	}

	if data.Count != 0 {
		return true
	}

	return false
}

func (c *Client) Close(ctx context.Context) error {
	return c.client.Close()
}
