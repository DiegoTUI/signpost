package utils_test

import (
	"math"
	"strconv"
	"testing"

	"fmt"

	"strings"

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

	random := utils.RandomInt(0, 0)
	if random != 0 {
		t.Error("RandomInt with two zeroes failed - ", random)
	}
}

func TestGetExternalIP(t *testing.T) {
	externalIp, err := utils.GetExternalIP()
	if err != nil {
		t.Error("GetExternalIp returned an error")
	}

	if strings.HasSuffix(externalIp, "\n") {
		t.Error("GetExternalIp ends with a carriage return")
	}

	octects := strings.Split(externalIp, ".")

	if len(octects) != 4 {
		t.Error("Invalid number of octects")
	}

	for _, octect := range octects {
		if octectInt, _ := strconv.ParseInt(octect, 10, 16); octectInt > 255 {
			t.Error("Invalid octect")
		}
	}
}
