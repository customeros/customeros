package data_fields

import (
	"time"

	neo4jmodel "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

// Nil fields wil be skipped from update
type MarkdownEventFields struct {
	AppSource      *string                    `json:"appSource,omitempty"`
	Source         *enum.Source               `json:"source,omitempty"`
	CreatedAt      *time.Time                 `json:"createdAt,omitempty"`
	OrganizationId *string                    `json:"organizationId,omitempty"`
	Content        *string                    `json:"content,omitempty"`
	ExternalSystem *neo4jmodel.ExternalSystem `json:"externalSystem,omitempty"`
}

func (fields MarkdownEventFields) ExternalSystemAvailable() bool {
	return fields.ExternalSystem != nil && fields.ExternalSystem.Available()
}

func (fields MarkdownEventFields) Type() string {
	return "MarkdownEventFields"
}
