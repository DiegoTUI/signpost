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

	// read capitals
	raw, err := ioutil.ReadFile(usr.HomeDir + "/resources/capitals.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	var capitals []models.Capital
	json.Unmarshal(raw, &capitals)

	log.Printf("number of capitals: %d\n", len(capitals))

	start := time.Now()
	// read world cities
	file, err := os.Open(usr.HomeDir + "/resources/worldcities.txt")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
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

	elapsed := time.Since(start)
	log.Printf("number of worldcities: %d - ellapsed: %s\n", len(worldCities), elapsed)

	// start building cities
	var cities []models.City

	for i := range capitals {
		worldCity := searchInWorldCities(worldCities, capitals[i].Capital)

		if worldCity != nil {
			city := models.City{
				Name:       worldCity.City,
				Country:    capitals[i].Country,
				Difficulty: 5,
				IsCapital:  true,
				Location:   models.NewGeoJSONPoint(worldCity.Latitude, worldCity.Longitude),
			}

			cities = append(cities, city)
		} else {
			log.Printf("Could not find match for %s", capitals[i].Capital)
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
