package models

// GeoJSONPoint defines a GeoJSON piont
type GeoJSONPoint struct {
	Type        string     `bson:"type" json:"-"`
	Coordinates [2]float64 `bson:"coordinates" json:"coordinates"`
}

// NewGeoJSONPoint creates a new GeoJSONPoint
func NewGeoJSONPoint(lat, lon float64) GeoJSONPoint {
	return GeoJSONPoint{
		Type:        "Point",
		Coordinates: [2]float64{lon, lat},
	}
}
