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

	capital, err := models.NewCapital(jsonString)

	if err != nil {
		t.Error("Creating a valid capital returned an error")
	}

	if capital.Country != "Spain" {
		t.Error("Country parsed incorrectly")
	}

	if capital.Capital != "Madrid" {
		t.Error("Capital parsed incorrectly")
	}

	if capital.Type != "countryCapital" {
		t.Error("Type parsed incorrectly")
	}
}
