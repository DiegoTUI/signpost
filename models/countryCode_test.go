package models_test

import (
	"testing"

	"github.com/DiegoTUI/signpost/models"
)

func TestNewCountryCode(t *testing.T) {
	jsonString := []byte(`{
		"name": "Afghanistan",
		"alpha-2": "AF",
		"country-code": "004"
	}`)

	countryCode, err := models.NewCountryCode(jsonString)

	if err != nil {
		t.Error("Creating a valid country code returned an error")
	}

	if countryCode.Country != "Afghanistan" {
		t.Error("Country parsed incorrectly")
	}

	if countryCode.ISOCode != "af" {
		t.Error("ISOCode parsed incorrectly")
	}
}
