package data_fields

import (
	"time"
)

type JobRoleFields struct {
	AppSource   *string    `json:"appSource,omitempty"`
	Source      *string    `json:"source,omitempty"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	EndedAt     *time.Time `json:"endedAt,omitempty"`
	JobTitle    *string    `json:"jobTitle,omitempty"`
	Primary     *bool      `json:"primary,omitempty"`
	Description *string    `json:"description,omitempty"`
	Company     *string    `json:"company,omitempty"`
}
