package models

import (
	"github.com/DiegoTUI/signpost/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// City defines a city mofor signpostingdel
type City struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"cityId"`
	Name       string        `bson:"name" json:"name"`
	Country    string        `bson:"country" json:"country"`
	Difficulty uint8         `bson:"difficulty" json:"difficulty"`
	IsCapital  bool          `bson:"isCapital" json:"isCapital"`
	Location   GeoJSONPoint  `bson:"location" json:"location"`
}

// Collection returns the name of the collection for the MongoObject
func (c City) Collection() string {
	return "cities"
}

// EnsureIndexes ensures tghe indexes of a certain model
func (c City) EnsureIndexes() error {
	indexes := []mgo.Index{
		mgo.Index{
			Key: []string{"country"},
		},
		mgo.Index{
			Key: []string{"name"},
		},
	}
	return db.EnsureIndexes(c, indexes)
}

// Insert inserts a document in the DB
func (c City) Insert() error {
	return db.Insert(c)
}

// Upsert upserts a document in the DB
func (c City) Upsert() error {
	return db.Upsert(c, bson.M{"name": c.Name, "country": c.Country})
}
