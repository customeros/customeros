package data_fields

import (
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"time"
)

// Nil fields wil be skipped from update
type IssueFields struct {
	AppSource                 *string                    `json:"appSource,omitempty"`
	Source                    *string                    `json:"source,omitempty"`
	CreatedAt                 *time.Time                 `json:"createdAt,omitempty"`
	GroupId                   *string                    `json:"groupId,omitempty"`
	Subject                   *string                    `json:"subject,omitempty"`
	Description               *string                    `json:"description,omitempty"`
	Status                    *string                    `json:"status,omitempty"`
	Priority                  *string                    `json:"priority,omitempty"`
	ReportedByOrganizationId  *string                    `json:"reportedByOrganizationId,omitempty"`
	SubmittedByOrganizationId *string                    `json:"submittedByOrganizationId,omitempty"`
	SubmittedByUserId         *string                    `json:"submittedByUserId,omitempty"`
	ExternalSystem            *neo4jmodel.ExternalSystem `json:"externalSystem,omitempty"`
}

func (fields IssueFields) ExternalSystemAvailable() bool {
	return fields.ExternalSystem != nil && fields.ExternalSystem.Available()
}
