package models

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoModel is an interface for all the models to be persisted in mongo.
type MongoModel interface {
	// Collection returns the name of the collection for the MongoObject
	Collection() string

	// Indexes returns the indexes of the collection for the MongoObject
	Indexes() []mgo.Index

	// FindOneQuery returns the query to perform upserts and findOnes
	FindOneQuery() bson.M
}
