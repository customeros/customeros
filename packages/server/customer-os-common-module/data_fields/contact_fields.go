package data_fields

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"time"
)

type ContactFields struct {
	AppSource       *string               `json:"appSource,omitempty"`
	Source          *string               `json:"source,omitempty"`
	ExternalSystem  *model.ExternalSystem `json:"externalSystem,omitempty"`
	FirstName       *string               `json:"firstName,omitempty"`
	LastName        *string               `json:"lastName,omitempty"`
	Name            *string               `json:"name,omitempty"`
	Username        *string               `json:"username,omitempty"`
	ProfilePhotoUrl *string               `json:"profilePhotoUrl,omitempty"`
	Timezone        *string               `json:"timezone,omitempty"`
	Description     *string               `json:"description,omitempty"`
	Prefix          *string               `json:"prefix,omitempty"`
	CreatedAt       *time.Time            `json:"createdAt,omitempty"`
	Hide            *bool                 `json:"hide,omitempty"`
}

func (fields ContactFields) ExternalSystemAvailable() bool {
	return fields.ExternalSystem != nil && fields.ExternalSystem.Available()
}
