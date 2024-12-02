package data_fields

import (
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
)

type Integration string

type FlowEvent string

type FlowEventObject[T any] struct {
	Tenant      *string     `json:"tenant"`
	Event       FlowEvent   `json:"event"`
	Integration Integration `json:"integration"`
	Timestamp   *time.Time  `json:"timestamp"`
	Payload     *T          `json:"payload"`
}

type FlowEventDefinition struct {
	Event       FlowEvent
	Integration Integration
	Schema      interface{} // The defined event type
	Description string
}

type MeetingSummaryCreatedEvent struct {
	Source            enum.ExternalSystemId `json:externalSystem`
	ParticipantEmails *[]string             `json:"participantEmails"`
	Content           *string               `json:"content,omitempty"`
}
