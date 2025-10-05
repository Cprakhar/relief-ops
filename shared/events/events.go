package events

import "github.com/cprakhar/relief-ops/shared/types"

// Event types
const (
	ResourceCommandFind   = "resource.cmd.find"
	UserNotifyAdminReview = "user.notify.admin_review"
)

type DisasterEventCreatedPayload struct {
	DisasterID  string            `json:"disaster_id"`
	Location    types.Coordinates `json:"location"`
	Range       int               `json:"range"`
	VolunteerID string            `json:"volunteer_id"`
}
