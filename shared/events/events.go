package events

import "github.com/cprakhar/relief-ops/shared/types"

// Event types
const (
	ResourceCommandFind   = "resource.cmd.find"
	DisasterCommandDelete = "disaster.cmd.delete"
)

type DisasterEventCreatedPayload struct {
	DisasterID string            `json:"disaster_id"`
	Location   types.Coordinates `json:"location"`
}
