package dto

import neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"

type UpdateCustomFieldTemplate struct {
	Name        *string   `json:"name,omitempty"`
	Type        *string   `json:"type,omitempty"`
	ValidValues *[]string `json:"validValues,omitempty"`
	Order       *int64    `json:"order,omitempty"`
	Required    *bool     `json:"required,omitempty"`
	Length      *int64    `json:"length,omitempty"`
	Min         *int64    `json:"min,omitempty"`
	Max         *int64    `json:"max,omitempty"`
}

func New_UpdateCustomFieldTemplate_From_CustomFieldTemplateSaveFields(data neo4jrepository.CustomFieldTemplateSaveFields) UpdateCustomFieldTemplate {
	updateEvent := UpdateCustomFieldTemplate{}
	if data.UpdateName {
		updateEvent.Name = &data.Name
	}
	if data.UpdateType {
		updateEvent.Type = &data.Type
	}
	if data.UpdateValidValues {
		updateEvent.ValidValues = &data.ValidValues
	}
	if data.UpdateOrder {
		updateEvent.Order = data.Order
	}
	if data.UpdateRequired {
		updateEvent.Required = data.Required
	}
	if data.UpdateLength {
		updateEvent.Length = data.Length
	}
	if data.UpdateMin {
		updateEvent.Min = data.Min
	}
	if data.UpdateMax {
		updateEvent.Max = data.Max
	}
	return updateEvent
}
