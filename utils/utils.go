package utils

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// PrettyPrint pretty prints an object
func PrettyPrint(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// AngleToTheNorth calculates the angle of a route given by the origin-destination vector
// with respect to the North
func AngleToTheNorth(origin, destination s2.LatLng) s1.Angle {
	lonDiff := (destination.Lng - origin.Lng).Radians()
	y := math.Sin(lonDiff) * math.Cos(destination.Lat.Radians())
	x := math.Cos(origin.Lat.Radians())*math.Sin(destination.Lat.Radians()) -
		math.Sin(origin.Lat.Radians())*math.Cos(destination.Lat.Radians())*math.Cos(lonDiff)

	return s1.Angle(math.Atan2(y, x))
}

// EarthDistance calculates the distance between two points on earth
func EarthDistance(origin, destination s2.LatLng) float64 {
	const RADIUS = 6371000
	angle := destination.Distance(origin).Normalized().Radians()
	if angle > math.Pi {
		angle = 2*math.Pi - angle
	}
	return RADIUS * angle
}

// RandomInt returns a random int between min and max
func RandomInt(min, max int) int {
	if min == max {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// GetExternalIP retrieves the external IP of the machine running the server
func GetExternalIP() (string, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := netClient.Get("http://ipinfo.io/ip")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyString := string(body)

	return bodyString[:len(bodyString)-1], nil
}
