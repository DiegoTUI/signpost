package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/DiegoTUI/signpost/models"

	"errors"
)

var session *mgo.Session
var database *mgo.Database

// Connect establishes a connection with the provided dbhost
func Connect(dbhost, dbname string) error {
	var err error
	session, err = mgo.Dial(dbhost)
	if err != nil {
		return err
	}

	database = session.DB(dbname)
	return nil
}

// EnsureIndex adds the corresponding indexes for a MongoModel
func EnsureIndex(item models.MongoModel) error {
	if database == nil {
		return errors.New("Could not ensure indexes. Database not connected")
	}

	for _, index := range item.Indexes() {
		err := database.C(item.Collection()).EnsureIndex(index)
		if err != nil {
			return err
		}
	}

	return nil
}

// Insert inserts an element in the DB
func Insert(item models.MongoModel) error {
	if database == nil {
		return errors.New("Could not insert. Database not connected")
	}

	err := database.C(item.Collection()).Insert(&item)

	return err
}

// Upsert upserts an element in the DB using the primary key provided
func Upsert(item models.MongoModel) (err error) {
	_, err = database.C(item.Collection()).Upsert(item.FindOneQuery(), bson.M{"$set": item})
	return
}
