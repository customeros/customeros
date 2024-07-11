package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"time"
)

type LogEntryDataFields struct {
	Content              string
	ContentType          string
	StartedAt            *time.Time
	AuthorUserId         *string
	LoggedOrganizationId *string
}

type LogEntryFields struct {
	ID                 string
	Tenant             string
	LogEntryDataFields LogEntryDataFields
	Source             events.Source
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}
