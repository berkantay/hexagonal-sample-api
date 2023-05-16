package domain

type Coordinate struct {
	Latitude  float64 `json:"latitude,omitempty" form:"latitude"`
	Longitude float64 `json:"longitude,omitempty" form:"longitude"`
}
