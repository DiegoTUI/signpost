package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"fmt"

	"flag"

	"github.com/DiegoTUI/signpost/models"
	"github.com/namsral/flag"
)

func main() {
	env := "development"
	flag.StringVar(&env, "env", env, "Environment: 'development' or 'production'")
	flag.Parse()

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
