package index

import (
	"testing"
)

func TestCreateKey(t *testing.T) {
	lat := 37.7749
	lon := -122.4194
	resolution := 13

	expectedKey := "8d283082800b3bf"

	result := CreateKey(lat, lon, resolution)

	if result != expectedKey {
		t.Errorf("Expected key: %s, but got: %s", expectedKey, result)
	}
}
