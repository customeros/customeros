package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEmailInputToEntity(input *model.EmailInput) *entity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := entity.EmailEntity{
		RawEmail:      input.Email,
		Label:         utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
		Primary:       utils.IfNotNilBool(input.Primary),
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     utils.IfNotNilString(input.AppSource),
	}
	if len(emailEntity.AppSource) == 0 {
		emailEntity.AppSource = constants.AppSourceCustomerOsApi
	}
	return &emailEntity
}

func MapEmailUpdateInputToEntity(input *model.EmailUpdateInput) *entity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := entity.EmailEntity{
		Id:            input.ID,
		Label:         utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
		Primary:       utils.IfNotNilBool(input.Primary),
		RawEmail:      utils.IfNotNilString(input.Email),
		SourceOfTruth: neo4jentity.DataSourceOpenline,
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
		Email:         utils.StringPtrFirstNonEmptyNillable(entity.Email, entity.RawEmail),
		RawEmail:      utils.StringPtrNillable(entity.RawEmail),
		Label:         &label,
		Primary:       entity.Primary,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		EmailValidationDetails: &model.EmailValidationDetails{
			Validated:      entity.Validated,
			IsReachable:    entity.IsReachable,
			IsValidSyntax:  entity.IsValidSyntax,
			CanConnectSMTP: entity.CanConnectSMTP,
			AcceptsMail:    entity.AcceptsMail,
			HasFullInbox:   entity.HasFullInbox,
			IsCatchAll:     entity.IsCatchAll,
			IsDeliverable:  entity.IsDeliverable,
			IsDisabled:     entity.IsDisabled,
			Error:          entity.Error,
		},
	}
}
