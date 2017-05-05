package main

import (
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

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	raw, err := ioutil.ReadFile(usr.HomeDir + "/resources/capitals.json")
	var capitals []models.Capital
	json.Unmarshal(raw, &capitals)

	fmt.Printf("number of capitals: %d\n", len(capitals))

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
