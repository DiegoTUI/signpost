package main

import (
	"log"
	"os"
	"os/user"
	"runtime"
	"time"

	"flag"

	"github.com/spf13/viper"

	"strings"

	"github.com/DiegoTUI/signpost/db"
	"github.com/DiegoTUI/signpost/models"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

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

	// read capitals
	capitals, err := readCapitals(usr.HomeDir + "/resources/capitals.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("number of capitals: %d\n", len(*capitals))

	// read country codes
	countryCodes, err := readCountryCodes(usr.HomeDir + "/resources/countryCodes.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("number of country codes: %d\n", len(*countryCodes))
	//prettyPrint(countryCodes)

	// read world cities
	worldCities := make(chan *models.WorldCity, 200)
	go readWorldCities(usr.HomeDir+"/resources/worldcities.txt", worldCities)

	// start building cities
	var cities []models.City

	for worldCity := range worldCities {
		countryCode := (*countryCodes)[worldCity.CountryCode]

		if countryCode == nil {
			continue
		}

		key := strings.ToLower(countryCode.Country) + "-" + strings.ToLower(worldCity.City)
		capital := (*capitals)[key]

		if capital == nil {
			key = strings.ToLower(countryCode.Country) + "-" + strings.ToLower(worldCity.AccentCity)
			capital = (*capitals)[key]
		}

		if capital == nil {
			continue
		}

		city := models.City{
			Name:       capital.Capital,
			Country:    capital.Country,
			Difficulty: 5,
			IsCapital:  true,
			Location:   models.NewGeoJSONPoint(worldCity.Latitude, worldCity.Longitude),
		}

		delete(*capitals, key)

		cities = append(cities, city)
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
		_, err := cities[i].Upsert()
		//log.Println("city - ", cities[i].Name, cities[i].Country)
		//prettyPrint(changeInfo)
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
