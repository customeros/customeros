package comment

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"time"
)

type Comment struct {
	ID               string                       `json:"id"`
	Tenant           string                       `json:"tenant"`
	Content          string                       `json:"content"`
	ContentType      string                       `json:"contentType,omitempty"`
	AuthorUserId     string                       `json:"authorUserId,omitempty"`
	CommentedIssueId string                       `json:"commentedIssueId,omitempty"`
	Source           events.Source                `json:"source"`
	ExternalSystems  []commonmodel.ExternalSystem `json:"externalSystem"`
	CreatedAt        time.Time                    `json:"createdAt,omitempty"`
	UpdatedAt        time.Time                    `json:"updatedAt,omitempty"`
}

type CommentDataFields struct {
	Content          string
	ContentType      string
	AuthorUserId     *string
	CommentedIssueId *string
}
