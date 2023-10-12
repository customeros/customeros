package models

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"time"
)

type LogEntry struct {
	ID                    string                  `json:"id"`
	Tenant                string                  `json:"tenant"`
	Content               string                  `json:"content"`
	ContentType           string                  `json:"contentType,omitempty"`
	AuthorUserId          string                  `json:"authorUserId,omitempty"`
	LoggedOrganizationIds []string                `json:"loggedOrganizationIds,omitempty"`
	Source                cmnmod.Source           `json:"source"`
	ExternalSystems       []cmnmod.ExternalSystem `json:"externalSystem"`
	CreatedAt             time.Time               `json:"createdAt,omitempty"`
	UpdatedAt             time.Time               `json:"updatedAt,omitempty"`
	StartedAt             time.Time               `json:"startedAt,omitempty"`
	TagIds                []string                `json:"tagIds,omitempty"`
}
