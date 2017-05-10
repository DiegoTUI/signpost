package models

import (
	"github.com/DiegoTUI/signpost/db"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Signpost defines a signpost model
type Signpost struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"cityId"`
	Center     City          `bson:"center" json:"center"`
	Signs      []sign        `bson:"sign" json:"sign"`
	Difficulty uint8         `bson:"difficulty" json:"difficulty"`
}

type sign struct {
	City     City    `bson:"city" json:"city"`
	Angle    float64 `bson:"angle" json:"angle"`
	Distance float64 `bson:"distance" json:"distance"`
}

// Collection returns the name of the collection for the MongoObject
func (s Signpost) Collection() string {
	return "signposts"
}

// EnsureIndexes ensures tghe indexes of a certain model
func (s Signpost) EnsureIndexes() error {
	indexes := []mgo.Index{
		mgo.Index{
			Key: []string{"center"},
		},
	}
	return db.EnsureIndexes(s, indexes)
}

// Insert inserts a document in the DB
func (s Signpost) Insert() error {
	return db.Insert(s)
}

// Upsert upserts a document in the DB
func (s Signpost) Upsert() (*mgo.ChangeInfo, error) {
	return db.Upsert(s, bson.M{})
}
