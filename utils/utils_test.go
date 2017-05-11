package utils_test

import (
	"math"
	"testing"

	"fmt"

	"github.com/DiegoTUI/signpost/utils"
	"github.com/golang/geo/s2"
)

const (
	TOLERANCE = 0.000001
)

type testDuple struct {
	name        string
	origin      s2.LatLng
	destination s2.LatLng
	expected    float64
	tolerance   float64
}

func (t testDuple) String() string {
	return fmt.Sprintf("%s - %f", t.name, t.expected)
}

func TestAngleToTheNorth(t *testing.T) {
	zero := s2.LatLng{
		Lat: 0,
		Lng: 0,
	}

	north := s2.LatLng{
		Lat: math.Pi / 3,
		Lng: 0,
	}

	south := s2.LatLng{
		Lat: -math.Pi / 3,
		Lng: 0,
	}

	east := s2.LatLng{
		Lat: 0,
		Lng: math.Pi / 3,
	}

	west := s2.LatLng{
		Lat: 0,
		Lng: -math.Pi / 3,
	}

	tests := []testDuple{
		testDuple{
			name:        "zero-north",
			origin:      zero,
			destination: north,
			expected:    0,
			tolerance:   TOLERANCE,
		},
		testDuple{
			name:        "zero-east",
			origin:      zero,
			destination: east,
			expected:    90,
			tolerance:   TOLERANCE,
		},
		testDuple{
			name:        "zero-south",
			origin:      zero,
			destination: south,
			expected:    180,
			tolerance:   TOLERANCE,
		},
		testDuple{
			name:        "north-south",
			origin:      north,
			destination: south,
			expected:    180,
			tolerance:   TOLERANCE,
		},
		testDuple{
			name:        "zero-west",
			origin:      zero,
			destination: west,
			expected:    -90,
			tolerance:   TOLERANCE,
		},
		testDuple{
			name:        "south-east",
			origin:      south,
			destination: east,
			expected:    63.4349488,
			tolerance:   TOLERANCE,
		},
		testDuple{
			name:        "east-south",
			origin:      east,
			destination: south,
			expected:    -153.4349488,
			tolerance:   TOLERANCE,
		},
	}

	for i := range tests {
		calculated := utils.AngleToTheNorth(tests[i].origin, tests[i].destination)
		if math.Abs(calculated.Degrees()-tests[i].expected) > tests[i].tolerance {
			t.Error("AngleToTheNorth failed ", tests[i], calculated)
		}
	}
}
