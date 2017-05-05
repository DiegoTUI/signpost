package models_test

import (
	"testing"

	"github.com/DiegoTUI/signpost/models"
)

func TestNewCapital(t *testing.T) {
	jsonString := []byte(`{
		"country": "Spain",
		"capital":"Madrid",
		"type": "countryCapital"
	}`)

	city, err := models.NewCapital(jsonString)

	if err != nil {
		t.Error("Creating a valid captal returned an error")
	}

	if city.Country != "Spain" {
		t.Error("Country parsed incorrectly")
	}

	if city.Capital != "Madrid" {
		t.Error("Capital parsed incorrectly")
	}

	if city.Type != "countryCapital" {
		t.Error("Type parsed incorrectly")
	}
}
