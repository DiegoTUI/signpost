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
func Upsert(item interfaces.MongoInterface, findOneQuery bson.M) (*mgo.ChangeInfo, error) {
	return database.C(item.Collection()).Upsert(findOneQuery, bson.M{"$set": item})
}

// Find finds all the documents matching the given query
func Find(query bson.M, collection string, results interface{}) error {
	return database.C(collection).Find(query).All(results)
}

// FindOne finds one document matching the given query
func FindOne(query bson.M, result interface{}) error {
	mongoResult, ok := result.(interfaces.MongoInterface)
	if !ok {
		return errors.New("Result is not a mongo interface")
	}
	return database.C(mongoResult.Collection()).Find(query).One(result)
}

// GetDB returns the current DB
func GetDB() *mgo.Database {
	return database
}
