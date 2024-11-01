package dto

import (
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"time"
)

type CreateContact struct {
	FirstName       string                     `json:"firstName"`
	LastName        string                     `json:"lastName"`
	Prefix          string                     `json:"prefix"`
	Description     string                     `json:"description"`
	Timezone        string                     `json:"timezone"`
	ProfilePhotoUrl string                     `json:"profilePhotoUrl"`
	Username        string                     `json:"username"`
	Name            string                     `json:"name"`
	Source          string                     `json:"sourceFields"`
	CreatedAt       time.Time                  `json:"createdAt"`
	ExternalSystem  *neo4jmodel.ExternalSystem `json:"externalSystem,omitempty"`
}

func New_CreateContact_From_ContactFields(data neo4jrepository.ContactFields, externalSystem neo4jmodel.ExternalSystem) CreateContact {
	output := CreateContact{
		FirstName:       data.FirstName,
		LastName:        data.LastName,
		Prefix:          data.Prefix,
		Description:     data.Description,
		Timezone:        data.Timezone,
		ProfilePhotoUrl: data.ProfilePhotoUrl,
		Username:        data.Username,
		Name:            data.Name,
		Source:          data.SourceFields.GetSource(),
		CreatedAt:       data.CreatedAt,
	}
	if externalSystem.Available() {
		output.ExternalSystem = &externalSystem
	}
	return output
}
