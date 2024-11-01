package dto

import neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"

type CreateCustomFieldTemplate struct {
	Name        string   `json:"name"`
	EntityType  string   `json:"entityType"`
	Type        string   `json:"type"`
	ValidValues []string `json:"validValues"`
	Order       *int64   `json:"order,omitempty"`
	Required    *bool    `json:"required,omitempty"`
	Length      *int64   `json:"length,omitempty"`
	Min         *int64   `json:"min,omitempty"`
	Max         *int64   `json:"max,omitempty"`
}

func New_CreateCustomFieldTemplate_From_CustomFieldTemplateSaveFields(data neo4jrepository.CustomFieldTemplateSaveFields) CreateCustomFieldTemplate {
	return CreateCustomFieldTemplate{
		Name:        data.Name,
		EntityType:  data.EntityType.String(),
		Type:        data.Type,
		ValidValues: data.ValidValues,
		Order:       data.Order,
		Required:    data.Required,
		Length:      data.Length,
		Min:         data.Min,
		Max:         data.Max,
	}
}
