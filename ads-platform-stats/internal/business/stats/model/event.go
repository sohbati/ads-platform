package model

import "time"

const (
	EventView           = "view"
	EventContactReveal  = "contact_reveal"
	EventCall           = "call"
)

type Event struct {
	AdID          int64     `json:"ad_id"`
	Event         string    `json:"event"`
	ViewerID      string    `json:"viewer_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	SessionUserID *int64    `json:"session_user_id,omitempty"`
}

func ValidEvent(name string) bool {
	switch name {
	case EventView, EventContactReveal, EventCall:
		return true
	default:
		return false
	}
}
