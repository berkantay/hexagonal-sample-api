package tiles

import (
	"context"

	"github.com/berkantay/firefly-weather-condition-api/config"
	"github.com/berkantay/firefly-weather-condition-api/internal/domain"
	"github.com/xjem/t38c"
)

type Client struct {
	client *t38c.Client
}

func NewClient(config *config.Config) (*Client, error) {
	client, err := t38c.New(t38c.Config{
		Address: "localhost:9851",
		Debug:   false, // print queries to stdout
	})
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		client: client,
	}, nil
}

func (c *Client) CityIntersectByCode(ctx context.Context, code string, coordinate *domain.Coordinate) bool {
	data, err := c.client.Search.Intersects("cities").Circle(coordinate.Latitude, coordinate.Longitude, 0).Do(ctx)
	if err != nil {
		return false
	}

	for _, c := range data.Objects {
		if code == c.ID {
			return true
		}
	}
	return false
}
