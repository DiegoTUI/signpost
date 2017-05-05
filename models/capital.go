package models

import (
	"encoding/json"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Capital defines a capital model
type Capital struct {
	Country string `json:"country"`
	Capital string `json:"capital"`
	Type    string `json:"type"`
}

// NewCapital creates a new capital from a JSON
func NewCapital(jsonByte []byte) (*Capital, error) {
	var capital Capital

	if err := json.Unmarshal(jsonByte, &capital); err != nil {
		return nil, err
	}

	return &capital, nil
}

// Collection returns the name of the collection for the MongoObject
func (c Capital) Collection() string {
	return "capitals"
}

// Indexes returns the indexes for the collection in mongoDB
func (c Capital) Indexes() []mgo.Index {
	return []mgo.Index{
		mgo.Index{
			Key: []string{"country"},
		},
		mgo.Index{
			Key:    []string{"capital"},
			Unique: true,
		},
	}
}

// FindOneQuery returns the main query to perform upserts and findOnes
func (c Capital) FindOneQuery() bson.M {
	return bson.M{"capital": c.Capital}
}
