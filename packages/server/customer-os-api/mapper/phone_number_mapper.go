package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapPhoneNumberInputToEntity(input *model.PhoneNumberInput) *neo4jentity.PhoneNumberEntity {
	if input == nil {
		return nil
	}
	phoneNumberEntity := neo4jentity.PhoneNumberEntity{
		RawPhoneNumber: input.PhoneNumber,
		Label:          utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
		Primary:        utils.IfNotNilBool(input.Primary),
		Source:         neo4jentity.DataSourceOpenline,
		SourceOfTruth:  neo4jentity.DataSourceOpenline,
	}
	return &phoneNumberEntity
}

func MapEntitiesToPhoneNumbers(entities *neo4jentity.PhoneNumberEntities) []*model.PhoneNumber {
	var phoneNumbers []*model.PhoneNumber
	for _, phoneNumberEntity := range *entities {
		phoneNumbers = append(phoneNumbers, MapEntityToPhoneNumber(&phoneNumberEntity))
	}
	return phoneNumbers
}

func MapEntityToPhoneNumber(entity *neo4jentity.PhoneNumberEntity) *model.PhoneNumber {
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
