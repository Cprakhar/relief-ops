package types

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	Tags        []string      `json:"tags" bson:"tags"`
	VolunteerID string        `json:"volunteer_id" bson:"volunteer_id"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
	ImageURLs   []string      `json:"image_urls" bson:"image_urls"`
	Location    Coordinates   `json:"location" bson:"location"`
	Status      string        `json:"status" bson:"status"`
}

type User struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	Email     string        `json:"email" bson:"email"`
	Password  string        `json:"-" bson:"password"`
	AvatarURL string        `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
	Role      string        `json:"role" bson:"role"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

type Resource struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	AmenityType string        `json:"amenity_type" bson:"amenity_type"` // e.g., amentiy type
	Location    *Location     `json:"location" bson:"location"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

type Location struct {
	Type        string    `json:"type" bson:"type"`               // e.g., "Point"
	Coordinates []float64 `json:"coordinates" bson:"coordinates"` // [longitude, latitude]
}
