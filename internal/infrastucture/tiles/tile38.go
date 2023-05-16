package tiles

import (
	"github.com/berkantay/firefly-weather-condition-api/config"
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

func (c *Client) Client() *t38c.Client {
	return c.client
}
