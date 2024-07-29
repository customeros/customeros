package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
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
	Source             common.Source
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}
