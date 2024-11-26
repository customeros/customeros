// @openapi 3.0.0
package events

import "time"

type Event struct {
	ID              string
	Content         string
	EventTimestamp  time.Time
	Tenant          string
	OrganizationIDs []string
	Source          string
	SourceLogoUrl   string
}

type EventSource string

const (
	EventFathom EventSource = "FATHOM"
	EventGrain  EventSource = "GRAIN"
)
