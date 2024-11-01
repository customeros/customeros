package dto

import (
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"time"
)

type UpdateContact struct {
	FirstName         *string                    `json:"firstName,omitempty"`
	LastName          *string                    `json:"lastName,omitempty"`
	Prefix            *string                    `json:"prefix,omitempty"`
	Description       *string                    `json:"description,omitempty"`
	Timezone          *string                    `json:"timezone,omitempty"`
	ProfilePhotoUrl   *string                    `json:"profilePhotoUrl,omitempty"`
	Username          *string                    `json:"username,omitempty"`
	Name              *string                    `json:"name,omitempty"`
	Source            *string                    `json:"sourceFields,omitempty"`
	CreatedAt         *time.Time                 `json:"createdAt,omitempty"`
	UpdateOnlyIfEmpty bool                       `json:"updateOnlyIfEmpty,omitempty"`
	ExternalSystem    *neo4jmodel.ExternalSystem `json:"externalSystem,omitempty"`
}

func New_UpdateContact_From_ContactFields(data neo4jrepository.ContactFields, externalSystem neo4jmodel.ExternalSystem) UpdateContact {
	output := UpdateContact{
		UpdateOnlyIfEmpty: data.UpdateOnlyIfEmpty,
	}

	if data.UpdateFirstName {
		output.FirstName = &data.FirstName
	}
	if data.UpdateLastName {
		output.LastName = &data.LastName
	}
	if data.UpdatePrefix {
		output.Prefix = &data.Prefix
	}
	if data.UpdateDescription {
		output.Description = &data.Description
	}
	if data.UpdateTimezone {
		output.Timezone = &data.Timezone
	}
	if data.UpdateProfilePhotoUrl {
		output.ProfilePhotoUrl = &data.ProfilePhotoUrl
	}
	if data.UpdateUsername {
		output.Username = &data.Username
	}
	if data.UpdateName {
		output.Name = &data.Name
	}

	if externalSystem.Available() {
		output.ExternalSystem = &externalSystem
	}
	return output
}
