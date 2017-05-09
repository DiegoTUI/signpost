package models

import (
	"errors"
	"strconv"
	"strings"
)

// WorldCity defines a World city model
type WorldCity struct {
	CountryCode string
	City        string
	AccentCity  string
	Latitude    float64
	Longitude   float64
}

// NewWorldCity creates a new world city from a CSV line
func NewWorldCity(line string) (*WorldCity, error) {
	fields := strings.Split(line, ",")
	if len(fields) != 7 {
		return nil, errors.New("Unable to parse line")
	}

	countryCode := fields[0]
	city := fields[1]
	accentCity := fields[2]

	var latitude, longitude float64
	var err error
	// assert types
	if latitude, err = strconv.ParseFloat(fields[5], 64); err != nil {
		return nil, errors.New("Unable to parse latitude")
	}
	if longitude, err = strconv.ParseFloat(fields[6], 64); err != nil {
		return nil, errors.New("Unable to parse longitude")
	}

	var worldCity = WorldCity{
		CountryCode: countryCode,
		City:        city,
		AccentCity:  accentCity,
		Latitude:    latitude,
		Longitude:   longitude,
	}

	return &worldCity, nil
}
