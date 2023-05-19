package index

import "github.com/uber/h3-go/v4"

// Create H3 key with given latitude, longitude and resolution
func CreateKey(lat, lon float64, resoultion int) string {
	coordinates := h3.NewLatLng(lat, lon)
	return h3.LatLngToCell(coordinates, resoultion).String()
}
