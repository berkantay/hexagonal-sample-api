package domain

type Coordinate struct {
	Latitude  float64 `json:"latitude,omitempty" form:"latitude"  binding:"required"` //these could be pointer because request may be nil but these are initialized as 0
	Longitude float64 `json:"longitude,omitempty" form:"longitude"  binding:"required"`
}
