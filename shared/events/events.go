package events

import "github.com/cprakhar/relief-ops/shared/types"

// Event types
const (
	ResourceCommandFind   = "resource.cmd.find"
	DisasterCommandDelete = "disaster.cmd.delete"
	UserNotifyAdminReview = "user.notify.admin_review"
)

type DisasterEventCreatedPayload struct {
	DisasterID    string            `json:"disaster_id"`
	Location      types.Coordinates `json:"location"`
	ContributorID string            `json:"contributor_id"`
}
