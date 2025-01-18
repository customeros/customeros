package data_fields

import (
	neo4jmodel "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"
	"time"
)

// Nil fields wil be skipped from update
type LogEntryFields struct {
	ID             string                     `json:"id,omitempty"`
	AppSource      *string                    `json:"appSource,omitempty"`
	Source         *string                    `json:"source,omitempty"`
	CreatedAt      *time.Time                 `json:"createdAt,omitempty"`
	StartedAt      *time.Time                 `json:"startedAt,omitempty" `
	OrganizationId *string                    `json:"organizationId,omitempty"`
	Content        *string                    `json:"content,omitempty"`
	ContentType    *string                    `json:"contentType,omitempty"`
	AuthorUserId   *string                    `json:"authorUserId,omitempty"`
	ExternalSystem *neo4jmodel.ExternalSystem `json:"externalSystem,omitempty"`
}

func (fields LogEntryFields) ExternalSystemAvailable() bool {
	return fields.ExternalSystem != nil && fields.ExternalSystem.Available()
}
