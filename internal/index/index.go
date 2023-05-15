package index

import "github.com/uber/h3-go/v4"

func CreatKey(lat, lon float64, resoultion int) string {
	coordinates := h3.NewLatLng(lat, lon)
	return h3.LatLngToCell(coordinates, resoultion).String()
}
