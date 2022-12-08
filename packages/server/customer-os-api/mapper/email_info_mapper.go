package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapEmailInputToEntity(input *model.EmailInput) *entity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := entity.EmailEntity{
		Email:   input.Email,
		Label:   utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
		Primary: utils.IfNotNilBool(input.Primary),
	}
	return &emailEntity
}

func MapEmailUpdateInputToEntity(input *model.EmailUpdateInput) *entity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := entity.EmailEntity{
		Id:      input.ID,
		Email:   input.Email,
		Label:   utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
		Primary: utils.IfNotNilBool(input.Primary),
	}
	return &emailEntity
}

func MapEntitiesToEmails(entities *entity.EmailEntities) []*model.Email {
	var emails []*model.Email
	for _, emailEntity := range *entities {
		emails = append(emails, MapEntityToEmail(&emailEntity))
	}
	return emails
}

func MapEntityToEmail(emailEntity *entity.EmailEntity) *model.Email {
	var label = model.EmailLabel(emailEntity.Label)
	if !label.IsValid() {
		label = ""
	}
	return &model.Email{
		ID:      emailEntity.Id,
		Email:   emailEntity.Email,
		Label:   &label,
		Primary: emailEntity.Primary,
	}
}
