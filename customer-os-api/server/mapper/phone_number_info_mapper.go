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
		Number: phoneNumberEntity.Number,
		Label:  label,
	}
}
