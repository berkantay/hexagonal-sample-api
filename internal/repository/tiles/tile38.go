package tiles

import (
	"context"
	"fmt"

	"github.com/berkantay/firefly-weather-condition-api/config"
	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/xjem/t38c"
)

type Client struct {
	client *t38c.Client
}

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
