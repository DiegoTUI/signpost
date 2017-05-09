package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/DiegoTUI/signpost/models"
)

func readCapitals(filePath string) (*map[string]*models.Capital, error) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var capitals []models.Capital
	json.Unmarshal(raw, &capitals)

	result := make(map[string]*models.Capital)
	for i := range capitals {
		result[strings.ToLower(capitals[i].Country)+"-"+strings.ToLower(capitals[i].Capital)] = &capitals[i]
	}

	return &result, nil
}

func readCountryCodes(filePath string) (*map[string]*models.CountryCode, error) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var countryCodes []models.CountryCode
	json.Unmarshal(raw, &countryCodes)

	result := make(map[string]*models.CountryCode)
	for i := range countryCodes {
		result[strings.ToLower(countryCodes[i].ISOCode)] = &countryCodes[i]
	}

	return &result, nil
}

func readWorldCities(filePath string, channel chan *models.WorldCity) {
	file, err := os.Open(filePath)
	if err != nil {
		close(channel)
	}
	scanner := bufio.NewScanner(file)

	// skip the first line
	scanner.Scan()

	// emit world cities as we read them
	for scanner.Scan() {
		line := scanner.Text()
		worldCity, err := models.NewWorldCity(line)
		if err != nil {
			log.Println("error parsing line: " + line)
		}
		channel <- worldCity
	}

	close(channel)
}

func searchInWorldCities(worldCities []models.WorldCity, city string) *models.WorldCity {
	for i := range worldCities {
		if strings.EqualFold(city, worldCities[i].City) ||
			strings.EqualFold(city, worldCities[i].AccentCity) {
			return &worldCities[i]
		}
	}

	return nil
}

func searchInCapitals(capitals *[]models.Capital, worldCity *models.WorldCity) *models.Capital {
	for i := range *capitals {
		if strings.EqualFold(worldCity.City, (*capitals)[i].Capital) ||
			strings.EqualFold(worldCity.AccentCity, (*capitals)[i].Capital) {
			result := &(*capitals)[i]
			*capitals = append((*capitals)[:i], (*capitals)[i+1:]...)
			return result
		}
	}

	return nil
}

func searchInCountryCodes(countryCodes []models.CountryCode, countryCode string) *models.CountryCode {
	for i := range countryCodes {
		if countryCodes[i].ISOCode == countryCode {
			return &countryCodes[i]
		}
	}

	return nil
}

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}
