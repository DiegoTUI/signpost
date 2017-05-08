package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/DiegoTUI/signpost/interfaces"

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

// Disconnect closes the current collection
func Disconnect() {
	session.Close()
}

// EnsureIndexes adds the corresponding indexes for a MongoModel
func EnsureIndexes(item interfaces.MongoInterface, indexes []mgo.Index) error {
	if database == nil {
		return errors.New("Could not ensure indexes. Database not connected")
	}

	for _, index := range indexes {
		err := database.C(item.Collection()).EnsureIndex(index)
		if err != nil {
			return err
		}
	}

	return nil
}

// Insert inserts an element in the DB
func Insert(item interfaces.MongoInterface) error {
	if database == nil {
		return errors.New("Could not insert. Database not connected")
	}

	err := database.C(item.Collection()).Insert(&item)

	return err
}

// Upsert upserts an element in the DB using the primary key provided
func Upsert(item interfaces.MongoInterface, findOneQuery bson.M) (err error) {
	_, err = database.C(item.Collection()).Upsert(findOneQuery, bson.M{"$set": item})
	return
}

// GetDB returns the current DB
func GetDB() *mgo.Database {
	return database
}
