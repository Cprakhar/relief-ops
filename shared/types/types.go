package types

import (
	"time"
)

const (
	Hospital    = "hospital"
	FireStation = "fire_station"
	Police      = "police"
	Shelter     = "shelter"
	Pharmacy    = "pharmacy"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Disaster struct {
	ID            string      `json:"id" bson:"_id,omitempty"`
	Title         string      `json:"title" bson:"title"`
	Description   string      `json:"description" bson:"description"`
	Tags          []string    `json:"tags" bson:"tags"`
	ContributorID string      `json:"contributor_id" bson:"contributor_id"`
	CreatedAt     time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" bson:"updated_at"`
	ImageURLs     []string    `json:"image_urls" bson:"image_urls"`
	Location      Coordinates `json:"location" bson:"location"`
	Status        string      `json:"status" bson:"status"`
}
