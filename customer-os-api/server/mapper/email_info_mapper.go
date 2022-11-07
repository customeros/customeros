package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapEmailInputToEntity(input *model.EmailInput) *entity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := entity.EmailEntity{
		Email: input.Email,
		Label: input.Label.String(),
	}
	if input.Primary != nil {
		emailEntity.Primary = *input.Primary
	} else {
		emailEntity.Primary = false
	}
	return &emailEntity
}

func MapEmailUpdateInputToEntity(input *model.EmailUpdateInput) *entity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := entity.EmailEntity{
		Id:    input.ID,
		Email: input.Email,
		Label: input.Label.String(),
	}
	if input.Primary != nil {
		emailEntity.Primary = *input.Primary
	} else {
		emailEntity.Primary = false
	}
	return &emailEntity
}

func MapEntitiesToEmails(entities *entity.EmailEntities) []*model.EmailInfo {
	var emails []*model.EmailInfo
	for _, emailEntity := range *entities {
		emails = append(emails, MapEntityToEmail(&emailEntity))
	}
	return emails
}

func MapEntityToEmail(emailEntity *entity.EmailEntity) *model.EmailInfo {
	var label = model.EmailLabel(emailEntity.Label)
	if !label.IsValid() {
		label = model.EmailLabelOther
	}
	return &model.EmailInfo{
		ID:      emailEntity.Id,
		Email:   emailEntity.Email,
		Label:   label,
		Primary: emailEntity.Primary,
	}
}
