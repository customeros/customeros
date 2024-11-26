package data_fields

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"time"
)

// Nil fields wil be skipped from update
type MarkdownEventFields struct {
	AppSource      *string                    `json:"appSource,omitempty"`
	Source         *neo4jentity.DataSource    `json:"source,omitempty"`
	CreatedAt      *time.Time                 `json:"createdAt,omitempty"`
	OrganizationId *string                    `json:"organizationId,omitempty"`
	Content        *string                    `json:"content,omitempty"`
	ExternalSystem *neo4jmodel.ExternalSystem `json:"externalSystem,omitempty"`
}

func (fields MarkdownEventFields) ExternalSystemAvailable() bool {
	return fields.ExternalSystem != nil && fields.ExternalSystem.Available()
}
