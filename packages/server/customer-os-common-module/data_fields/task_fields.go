package data_fields

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"time"
)

// Nil fields wil be skipped from update
type TaskFields struct {
	Subject         *string          `json:"subject,omitempty"`
	Description     *string          `json:"description,omitempty"`
	Status          *enum.TaskStatus `json:"status,omitempty"`
	DueAt           *time.Time       `json:"dueAt,omitempty"`
	OpportunityIds  *[]string        `json:"opportunityIds,omitempty"`
	AssigneeUserIds *[]string        `json:"assigneeUserIds,omitempty"`
	AppSource       *string          `json:"appSource,omitempty"`
	Source          *string          `json:"source,omitempty"`
	CreatedByUserId *string          `json:"reportedBy,omitempty"`
}
