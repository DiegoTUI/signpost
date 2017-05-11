package utils

import "encoding/json"
import "github.com/golang/geo/s2"
import "github.com/golang/geo/s1"
import "math"

// PrettyPrint pretty prints an object
func PrettyPrint(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// AngleToTheNorth calculates the angle of a route given by the origin-destination vector
// with respect to the North
func AngleToTheNorth(origin, destination s2.LatLng) s1.Angle {
	lonDiff := float64(destination.Lng - origin.Lng)
	y := math.Sin(lonDiff) * math.Cos(float64(destination.Lat))
	x := math.Cos(float64(origin.Lat))*math.Sin(float64(destination.Lat)) -
		math.Sin(float64(origin.Lat))*math.Cos(float64(destination.Lat))*math.Cos(lonDiff)

	return s1.Angle(math.Atan2(y, x))
}
