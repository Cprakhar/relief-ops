package types

import (
	"time"

	"github.com/cprakhar/relief-ops/shared/types"
)

type Disaster struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	Tags          []string          `json:"tags"`
	ContributorID string            `json:"contributor_id"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	ImageURLs     []string          `json:"image_urls"`
	Location      types.Coordinates `json:"location"`
	Status        string            `json:"status"`
}
