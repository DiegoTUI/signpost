package models

import (
	"log"

	"github.com/DiegoTUI/signpost/db"
	"github.com/DiegoTUI/signpost/utils"
	"github.com/golang/geo/s2"
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

// NewSignpost creates a new signpost with the given parameters
func NewSignpost(center City,
	minNumberOfSigns, maxNumberOfSigns uint8,
	minDistance, maxDistance float64,
	minDifficulty, maxDifficulty uint8) (*Signpost, error) {
	// check for the obvious
	if maxNumberOfSigns < minNumberOfSigns ||
		maxDistance < minDistance ||
		maxDifficulty < minDifficulty {
		return nil, nil
	}

	// build the geo query to the cities collection
	query := bson.M{
		"difficulty": bson.M{
			"$gte": minDifficulty,
			"$lte": maxDifficulty,
		},
		"location": bson.M{
			"$near": bson.M{
				"geometry": bson.M{
					"type":        "Point",
					"coordinates": center.Location.Coordinates,
				},
				"$maxDistance": maxDistance,
				"$minDistance": minDistance,
			},
		},
	}

	// perform the query
	var cities []City

	err := db.Find(query, City{}.Collection(), &cities)
	if err != nil {
		log.Println("error in find", err)
		return nil, err
	}

	// check for the minimum number of cities
	if len(cities) < int(minNumberOfSigns) {
		return nil, nil
	}

	// set the number of sities for the signpost
	numberOfCities := len(cities)
	if numberOfCities > int(maxNumberOfSigns) {
		numberOfCities = int(maxNumberOfSigns)
	}

	cityDistribution := make(map[float64][]*City)

	// distribute cities
	for i := range cities {
		city := &cities[i]
		latLngCenter := s2.LatLngFromDegrees(center.Location.Coordinates[1],
			center.Location.Coordinates[0])
		latLngCity := s2.LatLngFromDegrees(city.Location.Coordinates[1],
			city.Location.Coordinates[0])
		angle := utils.AngleToTheNorth(latLngCenter, latLngCity)

		log.Println(cityDistribution, angle)
	}

	return nil, nil
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
	query := bson.M{
		"center.name":    s.Center.Name,
		"center.country": s.Center.Country,
		"signs": bson.M{
			"$size": len(s.Signs),
			"$all":  s.Signs,
		},
	}
	return db.Upsert(s, query)
}
