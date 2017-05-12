package utils_test

import (
	"math"
	"testing"

	"fmt"

	"github.com/DiegoTUI/signpost/utils"
	"github.com/golang/geo/s2"
)

const (
	TOLERANCE = 0.001
)

type testDuple struct {
	name             string
	origin           s2.LatLng
	destination      s2.LatLng
	expectedAngle    float64
	expectedDistance float64
	tolerance        float64
}

func (t testDuple) String() string {
	return fmt.Sprintf("%s - %f - %f", t.name, t.expectedAngle, t.expectedDistance)
}

var zero = s2.LatLng{
	Lat: 0,
	Lng: 0,
}

var north = s2.LatLng{
	Lat: math.Pi / 3,
	Lng: 0,
}

var south = s2.LatLng{
	Lat: -math.Pi / 3,
	Lng: 0,
}

var east = s2.LatLng{
	Lat: 0,
	Lng: math.Pi / 3,
}

var west = s2.LatLng{
	Lat: 0,
	Lng: -math.Pi / 3,
}

var tests = []testDuple{
	testDuple{
		name:             "zero-north",
		origin:           zero,
		destination:      north,
		expectedAngle:    0,
		expectedDistance: 6672000,
		tolerance:        TOLERANCE,
	},
	testDuple{
		name:             "zero-east",
		origin:           zero,
		destination:      east,
		expectedAngle:    90,
		expectedDistance: 6672000,
		tolerance:        TOLERANCE,
	},
	testDuple{
		name:             "zero-south",
		origin:           zero,
		destination:      south,
		expectedAngle:    180,
		expectedDistance: 6672000,
		tolerance:        TOLERANCE,
	},
	testDuple{
		name:             "north-south",
		origin:           north,
		destination:      south,
		expectedAngle:    180,
		expectedDistance: 13343000,
		tolerance:        TOLERANCE,
	},
	testDuple{
		name:             "zero-west",
		origin:           zero,
		destination:      west,
		expectedAngle:    -90,
		expectedDistance: 6672000,
		tolerance:        TOLERANCE,
	},
	testDuple{
		name:             "south-east",
		origin:           south,
		destination:      east,
		expectedAngle:    63.4349488,
		expectedDistance: 8398000,
		tolerance:        TOLERANCE,
	},
	testDuple{
		name:             "east-south",
		origin:           east,
		destination:      south,
		expectedAngle:    -153.4349488,
		expectedDistance: 8398000,
		tolerance:        TOLERANCE,
	},
}

func TestAngleToTheNorth(t *testing.T) {
	for i := range tests {
		calculated := utils.AngleToTheNorth(tests[i].origin, tests[i].destination)
		if math.Abs(calculated.Degrees()-tests[i].expectedAngle) > tests[i].tolerance {
			t.Error("AngleToTheNorth failed ", tests[i], calculated)
		}
	}
}

func TestEarthDistance(t *testing.T) {
	for i := range tests {
		calculated := utils.EarthDistance(tests[i].origin, tests[i].destination)
		if math.Abs(calculated-tests[i].expectedDistance) > 1000 {
			t.Error("EarthDistance failed ", tests[i], calculated)
		}
	}
}

func TestRandomInt(t *testing.T) {
	for i := 0; i < 10000; i++ {
		random := utils.RandomInt(0, 10)
		if random < 0 || random >= 10 {
			t.Error("RandomInt failed - ", random)
		}
	}
}
