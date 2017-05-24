package models

import (
	"errors"
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
	Signs      []Sign        `bson:"sign" json:"sign"`
	Difficulty uint8         `bson:"difficulty" json:"difficulty"`
}

// Sign defines a sign model
type Sign struct {
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
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": center.Location.Coordinates,
				},
				"$maxDistance": maxDistance,
				"$minDistance": minDistance + 15, // to prevent the center city to appear
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
		log.Println("not enough cities")
		return nil, nil
	}

	// create the signpost with no signs
	result := Signpost{
		Center:     center,
		Difficulty: 5,
	}

	// set the number of signs for the signpost
	numberOfSigns := len(cities)
	if numberOfSigns > int(maxNumberOfSigns) {
		numberOfSigns = int(maxNumberOfSigns)
	}

	signDistribution := make([][]*Sign, numberOfSigns, numberOfSigns)
	portion := 360 / numberOfSigns

	// distribute cities
	for i := range cities {
		city := &cities[i]
		latLngCenter := s2.LatLngFromDegrees(center.Location.Coordinates[1],
			center.Location.Coordinates[0])
		latLngCity := s2.LatLngFromDegrees(city.Location.Coordinates[1],
			city.Location.Coordinates[0])
		angle := utils.AngleToTheNorth(latLngCenter, latLngCity).Degrees()
		if angle < 0 {
			angle += 360
		}
		distance := utils.EarthDistance(latLngCenter, latLngCity)

		sign := Sign{
			City:     *city,
			Angle:    angle,
			Distance: distance,
		}

		index := int(angle) / portion
		signDistribution[index] = append(signDistribution[index], &sign)
	}

	// select signs
	for i := range signDistribution {
		offset := 0
		for {
			normalizedIndex := circularIndex(i+offset, len(signDistribution))
			sign, newFragment, err := selectSignFromFragment(signDistribution[normalizedIndex])
			if err != nil {
				offset = nextOffset(offset)
			} else {
				result.Signs = append(result.Signs, *sign)
				signDistribution[normalizedIndex] = newFragment
				break
			}
		}
	}

	return &result, nil
}

func circularIndex(index, length int) int {
	for index < 0 {
		index += length
	}

	return index % length
}

func nextOffset(currentOffset int) int {
	if currentOffset == 0 {
		return 1
	}
	if currentOffset > 0 {
		return -currentOffset
	}

	return -currentOffset + 1
}

// returns the sign, the new fragment without the sign or an error
func selectSignFromFragment(currentFragment []*Sign) (*Sign, []*Sign, error) {
	index := utils.RandomInt(0, len(currentFragment))
	sign, newFragment, err := SignArrayExtract(currentFragment, index)
	if err != nil {
		return nil, currentFragment, err
	}

	return sign, newFragment, nil
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

// SignArrayExtract receives a slice and extracts an element from it
// returns the element, the modified slice and an error
func SignArrayExtract(inputSlice []*Sign, index int) (*Sign, []*Sign, error) {
	if index >= len(inputSlice) {
		return nil, nil, errors.New("Index out of bounds")
	}

	cloned := append([]*Sign(nil), inputSlice...)
	element := cloned[index]
	newSlice := append(cloned[:index], cloned[index+1:]...)

	return element, newSlice, nil
}
