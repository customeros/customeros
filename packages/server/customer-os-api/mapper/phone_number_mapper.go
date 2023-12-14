package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapPhoneNumberInputToEntity(input *model.PhoneNumberInput) *entity.PhoneNumberEntity {
	if input == nil {
		return nil
	}
	phoneNumberEntity := entity.PhoneNumberEntity{
		RawPhoneNumber: input.PhoneNumber,
		Label:          utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
		Primary:        utils.IfNotNilBool(input.Primary),
		Source:         entity.DataSourceOpenline,
		SourceOfTruth:  entity.DataSourceOpenline,
	}
	return &phoneNumberEntity
}

func MapPhoneNumberUpdateInputToEntity(input *model.PhoneNumberUpdateInput) *entity.PhoneNumberEntity {
	if input == nil {
		return nil
	}
	phoneNumberEntity := entity.PhoneNumberEntity{
		Id:             input.ID,
		Label:          utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
		Primary:        utils.IfNotNilBool(input.Primary),
		RawPhoneNumber: utils.IfNotNilString(input.PhoneNumber),
		SourceOfTruth:  entity.DataSourceOpenline,
		Source:         entity.DataSourceOpenline,
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

func MapEntityToPhoneNumber(entity *entity.PhoneNumberEntity) *model.PhoneNumber {
	var label = model.PhoneNumberLabel(entity.Label)
	if !label.IsValid() {
		label = ""
	}
	return &model.PhoneNumber{
		ID:             entity.Id,
		E164:           utils.StringPtrNillable(entity.E164),
		RawPhoneNumber: utils.StringPtrNillable(entity.RawPhoneNumber),
		Validated:      entity.Validated,
		Label:          utils.ToPtr(label),
		Primary:        entity.Primary,
		Source:         MapDataSourceToModel(entity.Source),
		AppSource:      utils.StringPtrNillable(entity.AppSource),
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}
