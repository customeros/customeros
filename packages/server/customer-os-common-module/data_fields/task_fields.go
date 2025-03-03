package data_fields

import "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

// Nil fields wil be skipped from update
type TaskFields struct {
	AppSource     *string          `json:"appSource,omitempty"`
	Source        *string          `json:"source,omitempty"`
	Name          *string          `json:"name,omitempty"`
	Description   *string          `json:"description,omitempty"`
	Status        *enum.TaskStatus `json:"status,omitempty"`
	DueAt         *string          `json:"dueAt,omitempty"`
	Opportunities *[]string        `json:"opportunities,omitempty"`
	Assignees     *[]string        `json:"assignees,omitempty"`
}
