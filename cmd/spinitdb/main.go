package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"fmt"

	"flag"

	"github.com/spf13/viper"

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

	fmt.Printf("%s - %s\n", dbhost, dbname)

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

	fmt.Printf("number of capitals: %d\n", len(capitals))

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

	fmt.Printf("number of worldcities: %d\n", len(worldCities))

	// connect to the DB
	err = db.Connect(dbhost, dbname)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// add indexes
	err = db.EnsureIndex(models.Capital{})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for _, capital := range capitals {
		err = db.Upsert(capital)
		if err != nil {
			log.Println(err)
		}
	}
}
