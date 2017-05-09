package models

import (
	"encoding/json"
	"strings"
)

// CountryCode defines a country code model
type CountryCode struct {
	Country string `json:"name"`
	ISOCode string `json:"alpha-2"`
}

// NewCountryCode creates a new country code from a JSON
func NewCountryCode(jsonByte []byte) (*CountryCode, error) {
	var countryCode CountryCode

	if err := json.Unmarshal(jsonByte, &countryCode); err != nil {
		return nil, err
	}

	countryCode.ISOCode = strings.ToLower(countryCode.ISOCode)

	return &countryCode, nil
}
