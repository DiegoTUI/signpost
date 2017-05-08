package models_test

import (
	"testing"

	"github.com/DiegoTUI/signpost/models"
)

func TestNewGeoJSONPoint(t *testing.T) {
	lat, lon := 12.34564, 15.3782378

	geoJSONPoint := models.NewGeoJSONPoint(lat, lon)

	if geoJSONPoint.Type != "Point" {
		t.Error("GeoJson Type created incorectly")
	}

	if geoJSONPoint.Coordinates[0] != lon {
		t.Error("Longitude created incorrectly")
	}

	if geoJSONPoint.Coordinates[1] != lat {
		t.Error("Latitude created incorrectly")
	}
}
