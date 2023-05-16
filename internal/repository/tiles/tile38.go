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
		Debug:   true, // print queries to stdout
	})
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		client: client,
	}, nil
}

func (c *Client) CityIntersect(ctx context.Context, coordinate domain.Coordinate) bool {
	// data, err := c.client.Execute(ctx, "INTERSECTS", "cities POINT  40.731328 -74.067534")
	// if err != nil {
	// 	return false
	// }
	// fmt.Println(string(data))
	return true
}
