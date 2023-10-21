package models

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
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
	Source             cmnmod.Source
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}
