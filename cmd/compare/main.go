package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"

	"strings"

	"github.com/DiegoTUI/signpost/models"
)

func main() {
	// Read the home folder
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	start := time.Now()

	// read capitals
	capitals, err := readCapitals(usr.HomeDir + "/resources/capitals.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("number of capitals: %d\n", len(capitals))

	// read country codes
	countryCodes, err := readCountryCodes(usr.HomeDir + "/resources/countryCodes.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("number of country codes: %d\n", len(countryCodes))

	// read world cities
	worldCities := make(chan *models.WorldCity, 200)
	go readWorldCities(usr.HomeDir+"/resources/worldcities.txt", worldCities)

	var lastCountryCode string
	count := 0
	// check countries and country codes
	for worldCity := range worldCities {
		if worldCity.CountryCode == lastCountryCode {
			continue
		}

		lastCountryCode = worldCity.CountryCode

		countryCode := searchInCountryCodes(countryCodes, worldCity.CountryCode)
		if countryCode == nil {
			log.Printf("Country code not found for %s", worldCity.CountryCode)
			continue
		}

		capital := searchInCapitals(capitals, countryCode.Country)
		if capital == nil {
			log.Printf("Country not found in capitals for %s", countryCode.Country)
			continue
		}

		count++
	}

	ellapsed := time.Since(start)
	log.Printf("countries processed: %d - ellapsed: %s\n", count, ellapsed)
}

func readCapitals(filePath string) ([]models.Capital, error) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var capitals []models.Capital
	json.Unmarshal(raw, &capitals)

	return capitals, nil
}

func readCountryCodes(filePath string) ([]models.CountryCode, error) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var countryCodes []models.CountryCode
	json.Unmarshal(raw, &countryCodes)

	for i := range countryCodes {
		countryCodes[i].ISOCode = strings.ToLower(countryCodes[i].ISOCode)
	}

	return countryCodes, nil
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

func searchInCountryCodes(countryCodes []models.CountryCode, countryCode string) *models.CountryCode {
	for i := range countryCodes {
		if countryCodes[i].ISOCode == countryCode {
			return &countryCodes[i]
		}
	}

	return nil
}

func searchInCapitals(capitals []models.Capital, country string) *models.Capital {
	for i := range capitals {
		if capitals[i].Country == country {
			return &capitals[i]
		}
	}

	return nil
}
