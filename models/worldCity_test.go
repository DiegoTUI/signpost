package models_test

import (
	"testing"

	"github.com/DiegoTUI/signpost/models"
)

func TestNewWorldCity(t *testing.T) {
	line := "ad,andorra la vella,Andorra la Vella,07,20430,42.5,1.5166667"

	worldCity, err := models.NewWorldCity(line)

	if err != nil {
		t.Error("Creating a valid world city returned an error")
	}

	if worldCity.CountryCode != "ad" {
		t.Error("CountryCode parsed incorrectly")
	}

	if worldCity.City != "andorra la vella" {
		t.Error("City parsed incorrectly")
	}

	if worldCity.AccentCity != "Andorra la Vella" {
		t.Error("AccentCity parsed incorrectly")
	}

	if worldCity.Latitude != 42.5 {
		t.Error("Latitude parsed incorrectly")
	}

	if worldCity.Longitude != 1.5166667 {
		t.Error("Longitude parsed incorrectly")
	}
}
