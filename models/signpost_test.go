package models_test

import (
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/DiegoTUI/signpost/db"
	"github.com/DiegoTUI/signpost/models"
	"github.com/spf13/viper"
)

func TestArrayExtract(t *testing.T) {

	array := []*models.Sign{
		&models.Sign{
			Distance: 0,
		},
		&models.Sign{
			Distance: 1,
		},
		&models.Sign{
			Distance: 2,
		},
		&models.Sign{
			Distance: 3,
			Angle:    2,
		},
	}

	if _, _, err := models.SignArrayExtract(array, 4); err == nil {
		t.Error("ArrayExtract did not return index out of bounds")
	}

	// extract element from empty array
	emptyArray := []*models.Sign{}
	item, sliced, err := models.SignArrayExtract(emptyArray, 0)
	if err == nil {
		t.Error("Extracting an element from an empty array did NOT return an error")
	}

	// extract element in the middle
	item, sliced, err = models.SignArrayExtract(array, 2)
	if err != nil {
		t.Error("ArrayExtract did return an error for element 2")
	}
	if item.Distance != 2 {
		t.Error("ArrayExtract returned wrong item for element 2")
	}
	if len(sliced) != 3 ||
		sliced[0].Distance != 0 ||
		sliced[1].Distance != 1 ||
		sliced[2].Distance != 3 {
		t.Error("ArrayExtract returned wrong slice for element 2")
	}
	if len(array) != 4 ||
		array[0].Distance != 0 ||
		array[1].Distance != 1 ||
		array[2].Distance != 2 ||
		array[3].Distance != 3 {
		t.Error("ArrayExtract modified original array for element 2")
	}

	// extract last element
	item, sliced, err = models.SignArrayExtract(array, 3)
	if err != nil {
		t.Error("ArrayExtract did return an error for element 3")
	}
	if item.Distance != 3 {
		t.Error("ArrayExtract returned wrong item for element 3")
	}
	if len(sliced) != 3 ||
		sliced[0].Distance != 0 ||
		sliced[1].Distance != 1 ||
		sliced[2].Distance != 2 {
		t.Error("ArrayExtract returned wrong slice for element 3")
	}
	if len(array) != 4 ||
		array[0].Distance != 0 ||
		array[1].Distance != 1 ||
		array[2].Distance != 2 ||
		array[3].Distance != 3 {
		t.Error("ArrayExtract modified original array for element 3")
	}

	// extract first element
	item, sliced, err = models.SignArrayExtract(array, 0)
	if err != nil {
		t.Error("ArrayExtract did return an error for element 0")
	}
	if item.Distance != 0 {
		t.Error("ArrayExtract returned wrong item for element 0")
	}
	if len(sliced) != 3 ||
		sliced[0].Distance != 1 ||
		sliced[1].Distance != 2 ||
		sliced[2].Distance != 3 {
		t.Error("ArrayExtract returned wrong slice for element 0")
	}
	if len(array) != 4 ||
		array[0].Distance != 0 ||
		array[1].Distance != 1 ||
		array[2].Distance != 2 ||
		array[3].Distance != 3 {
		t.Error("ArrayExtract modified original array for element 0")
	}
}

func TestNewSignpost(t *testing.T) {
	// read config
	viper.SetConfigName("app")
	viper.AddConfigPath("../config")

	err := viper.ReadInConfig()
	if err != nil {
		t.Error("config file could not be read")
	}

	dbhost := viper.GetString("testing.dbhost")
	dbname := viper.GetString("testing.dbname")
	// connect to the db
	err = db.Connect(dbhost, dbname)
	if err != nil {
		t.Error("DB connection failed")
	}
	// get the city of Madrid
	center := models.City{}

	err = db.FindOne(bson.M{"name": "Madrid"}, &center)

	if err != nil {
		t.Error("findOne failed for Madrid")
	}

	// low radius
	signpost, err := models.NewSignpost(center, 3, 6, 0, 600000, 2, 7)

	if err != nil {
		t.Error("Creating a signpost failed", err)
	}

	if len(signpost.Signs) != 3 {
		t.Error("Creating a signpost for Madrid at 600km and min 3 failed", len(signpost.Signs))
	}

	checkCenterNotInSigns(t, center, signpost.Signs)

	// high radius
	signpost, err = models.NewSignpost(center, 3, 6, 0, 1000000, 2, 7)

	if err != nil {
		t.Error("Creating a signpost failed", err)
	}

	if len(signpost.Signs) != 6 {
		t.Error("Creating a signpost for Madrid at 1000km and max 6 failed", len(signpost.Signs))
	}

	checkCenterNotInSigns(t, center, signpost.Signs)
}

func checkCenterNotInSigns(t *testing.T, center models.City, signs []models.Sign) {
	for i := range signs {
		if signs[i].City == center {
			t.Error("Center was included in the signpost")
		}
	}
}
