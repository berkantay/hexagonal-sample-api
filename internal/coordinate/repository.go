package coordinate

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/xjem/t38c"
)

type GeospatialDatabase interface {
	Client() *t38c.Client
}

type CachePersistance interface {
	Client() *redis.Client
}

type WeatherClient interface {
	Client() *http.Client
}
