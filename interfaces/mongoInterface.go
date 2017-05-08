package interfaces

// MongoInterface is an interface for all the models to be persisted in mongo.
type MongoInterface interface {
	// Collection returns the name of the collection for the MongoObject
	Collection() string

	// EnsureIndexes ensures tghe indexes of a certain model
	EnsureIndexes() error

	// Insert inserts a document in the DB
	Insert() error

	// Upsert upserts a document in the DB
	Upsert() error
}
