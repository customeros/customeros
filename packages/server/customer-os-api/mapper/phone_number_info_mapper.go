package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapPhoneNumberInputToEntity(input *model.PhoneNumberInput) *entity.PhoneNumberEntity {
	if input == nil {
		return nil
	}
	phoneNumberEntity := entity.PhoneNumberEntity{
		E164:  input.E164,
		Label: input.Label.String(),
	}
	if input.Primary != nil {
		phoneNumberEntity.Primary = *input.Primary
	} else {
		phoneNumberEntity.Primary = false
	}
	return &phoneNumberEntity
}

func MapPhoneNumberUpdateInputToEntity(input *model.PhoneNumberUpdateInput) *entity.PhoneNumberEntity {
	if input == nil {
		return nil
	}
	phoneNumberEntity := entity.PhoneNumberEntity{
		Id:    input.ID,
		E164:  input.E164,
		Label: input.Label.String(),
	}
	if input.Primary != nil {
		phoneNumberEntity.Primary = *input.Primary
	} else {
		phoneNumberEntity.Primary = false
	}
	return &phoneNumberEntity
}

func MapEntitiesToPhoneNumbers(entities *entity.PhoneNumberEntities) []*model.PhoneNumber {
	var phoneNumbers []*model.PhoneNumber
	for _, phoneNumberEntity := range *entities {
		phoneNumbers = append(phoneNumbers, MapEntityToPhoneNumber(&phoneNumberEntity))
	}
	return phoneNumbers
}

func MapEntityToPhoneNumber(phoneNumberEntity *entity.PhoneNumberEntity) *model.PhoneNumber {
	var label = model.PhoneNumberLabel(phoneNumberEntity.Label)
	if !label.IsValid() {
		label = model.PhoneNumberLabelOther
	}
	return &model.PhoneNumber{
		ID:      phoneNumberEntity.Id,
		E164:    phoneNumberEntity.E164,
		Label:   label,
		Primary: phoneNumberEntity.Primary,
	}
}
