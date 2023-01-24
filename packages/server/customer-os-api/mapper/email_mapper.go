package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapEmailInputToEntity(input *model.EmailInput) *entity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := entity.EmailEntity{
		Email:         input.Email,
		Label:         utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
		Primary:       utils.IfNotNilBool(input.Primary),
		Source:        entity.DataSourceOpenline,
		SourceOfTruth: entity.DataSourceOpenline,
		AppSource:     utils.IfNotNilString(input.AppSource),
	}
	if len(emailEntity.AppSource) == 0 {
		emailEntity.AppSource = common.AppSourceCustomerOsApi
	}
	return &emailEntity
}

func MapEmailUpdateInputToEntity(input *model.EmailUpdateInput) *entity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := entity.EmailEntity{
		Id:            input.ID,
		Email:         input.Email,
		Label:         utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
		Primary:       utils.IfNotNilBool(input.Primary),
		SourceOfTruth: entity.DataSourceOpenline,
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

func MapEntityToEmail(entity *entity.EmailEntity) *model.Email {
	var label = model.EmailLabel(entity.Label)
	if !label.IsValid() {
		label = ""
	}
	return &model.Email{
		ID:            entity.Id,
		Email:         entity.Email,
		Label:         &label,
		Primary:       entity.Primary,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}
