package model

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type InteractionEvent struct {
	Tenant          string                  `json:"tenant"`
	ID              string                  `json:"id"`
	Source          cmnmod.Source           `json:"source"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
	Content         string                  `json:"content"`
	ContentType     string                  `json:"contentType"`
	Channel         string                  `json:"channel"`
	ChannelData     string                  `json:"channelData"`
	EventType       string                  `json:"eventType"`
	Identifier      string                  `json:"identifier"`
	Summary         string                  `json:"summary"`
	ActionItems     []string                `json:"actionItems"`
	ExternalSystems []cmnmod.ExternalSystem `json:"externalSystem"`
	PartOfIssueId   string                  `json:"partOfIssueId,omitempty"`
	PartOfSessionId string                  `json:"partOfSessionId,omitempty"`
}
