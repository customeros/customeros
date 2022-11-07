package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapPhoneNumberInputToEntity(input *model.PhoneNumberInput) *entity.PhoneNumberEntity {
	if input == nil {
		return nil
	}
	phoneNumberEntity := entity.PhoneNumberEntity{
		Number: input.Number,
		Label:  input.Label.String(),
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
		Id:     input.ID,
		Number: input.Number,
		Label:  input.Label.String(),
	}
	if input.Primary != nil {
		phoneNumberEntity.Primary = *input.Primary
	} else {
		phoneNumberEntity.Primary = false
	}
	return &phoneNumberEntity
}

func MapEntitiesToPhoneNumbers(entities *entity.PhoneNumberEntities) []*model.PhoneNumberInfo {
	var phoneNumbers []*model.PhoneNumberInfo
	for _, phoneNumberEntity := range *entities {
		phoneNumbers = append(phoneNumbers, MapEntityToPhoneNumber(&phoneNumberEntity))
	}
	return phoneNumbers
}

func MapEntityToPhoneNumber(phoneNumberEntity *entity.PhoneNumberEntity) *model.PhoneNumberInfo {
	var label = model.PhoneLabel(phoneNumberEntity.Label)
	if !label.IsValid() {
		label = model.PhoneLabelOther
	}
	return &model.PhoneNumberInfo{
		ID:      phoneNumberEntity.Id,
		Number:  phoneNumberEntity.Number,
		Label:   label,
		Primary: phoneNumberEntity.Primary,
	}
}
