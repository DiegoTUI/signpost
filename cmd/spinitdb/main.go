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

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	raw, err := ioutil.ReadFile(usr.HomeDir + "/resources/capitals.json")
	var capitals []models.Capital
	json.Unmarshal(raw, &capitals)

	fmt.Printf("number of capitals: %d\n", len(capitals))
}
