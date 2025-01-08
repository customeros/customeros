package data_fields

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"time"
)

type ActionFields struct {
	AppSource       *string                `json:"appSource,omitempty"`
	Source          *string                `json:"source,omitempty"`
	CreatedAt       *time.Time             `json:"createdAt,omitempty"`
	Content         *string                `json:"content,omitempty"`
	Metadata        *string                `json:"metadata,omitempty"`
	ActionType      *enum.ActionType       `json:"actionType,omitempty"`
	ExtraProperties map[string]interface{} `json:"extraProperties,omitempty"`
}
