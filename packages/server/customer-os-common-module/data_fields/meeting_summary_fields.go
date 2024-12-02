package data_fields

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
)

type MeetingSummaryFields struct {
	ExternalSystem    *model.ExternalSystem `json:"externalSystem,omitempty"`
	ParticipantEmails *[]string             `json:"participantEmails"`
	Content           *string               `json:"content,omitempty"`
	Timestamp         *time.Time            `json:"timestamp,omitempty"`
}

func (fields MeetingSummaryFields) ExternalSystemAvailable() bool {
	return fields.ExternalSystem != nil && fields.ExternalSystem.Available()
}
