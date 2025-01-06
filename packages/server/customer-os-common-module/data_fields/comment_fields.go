package data_fields

import (
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"time"
)

// Nil fields wil be skipped from update
type CommentFields struct {
	AppSource        *string                    `json:"appSource,omitempty"`
	Source           *string                    `json:"source,omitempty"`
	CreatedAt        *time.Time                 `json:"createdAt,omitempty"`
	ExternalSystem   *neo4jmodel.ExternalSystem `json:"externalSystem,omitempty"`
	CommentedIssueId *string                    `json:"commentedIssueId,omitempty"`
	Content          *string                    `json:"content,omitempty"`
	ContentType      *string                    `json:"contentType,omitempty"`
	AuthorUserId     *string                    `json:"authorUserId,omitempty"`
}

func (fields CommentFields) ExternalSystemAvailable() bool {
	return fields.ExternalSystem != nil && fields.ExternalSystem.Available()
}
