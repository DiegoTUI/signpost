package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"

	"flag"

	"github.com/spf13/viper"

	"strings"

	"github.com/DiegoTUI/signpost/db"
	"github.com/DiegoTUI/signpost/models"
)

func main() {
	// read environment
	var env string
	flag.StringVar(&env, "env", env, "Environment: 'development' or 'production'")
	flag.Parse()

	if env != "production" {
		env = "development"
	}

	// read config
	viper.SetConfigName("app")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file not found...")
		os.Exit(1)
	}

	dbhost := viper.GetString(env + ".dbhost")
	dbname := viper.GetString(env + ".dbname")

	log.Printf("%s - %s\n", dbhost, dbname)

	// Read the home folder
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	start := time.Now()

	// read world cities
	worldCities, err := readWorldCities(usr.HomeDir + "/resources/worldcities.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("number of worldcities: %d\n", len(worldCities))

	// read capitals
	capitals, err := readCapitals(usr.HomeDir + "/resources/capitals.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("number of capitals: %d\n", len(capitals))

	// start building cities
	var cities []models.City

	for i := range worldCities {
		capital := searchInCapitals(&capitals, &worldCities[i])

		if capital != nil {
			city := models.City{
				Name:       worldCities[i].City,
				Country:    capital.Country,
				Difficulty: 5,
				IsCapital:  true,
				Location:   models.NewGeoJSONPoint(worldCities[i].Latitude, worldCities[i].Longitude),
			}

			cities = append(cities, city)
		}
	}

	log.Printf("number of cities: %d\n", len(cities))

	// connect to the DB
	log.Println("Connecting to mongo")
	err = db.Connect(dbhost, dbname)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Println("saving cities")
	// save cities
	for i := range cities {
		err := cities[i].Upsert()
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("ensuring indexes")
	// ensure indexes
	err = models.City{}.EnsureIndexes()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Println("close db")
	// close DB
	db.Disconnect()

	ellapsed := time.Since(start)
	log.Printf("time ellapsed: %s\n", ellapsed)
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

func readWorldCities(filePath string) ([]models.WorldCity, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)

	// skip the first line
	scanner.Scan()

	// fill in an array of WorldCities
	var worldCities []models.WorldCity
	for scanner.Scan() {
		line := scanner.Text()
		worldCity, err := models.NewWorldCity(line)
		if err != nil {
			log.Println("error parsing line: " + line)
		}
		worldCities = append(worldCities, *worldCity)
	}

	return worldCities, nil
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
