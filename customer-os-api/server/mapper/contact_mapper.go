package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapContactInputToEntity(input model.ContactInput) *entity.ContactNode {
	return &entity.ContactNode{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Label:       *input.Label,
		ContactType: *input.ContactType,
	}
}

func MapEntityToContact(contact *entity.ContactNode) *model.Contact {
	return &model.Contact{
		ID:          contact.Id,
		FirstName:   contact.FirstName,
		LastName:    contact.LastName,
		Label:       &contact.Label,
		ContactType: &contact.ContactType,
	}
}
