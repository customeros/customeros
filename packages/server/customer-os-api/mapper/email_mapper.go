package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	enummapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEmailInputToEntity(input *model.EmailInput) *neo4jentity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := neo4jentity.EmailEntity{
		RawEmail: input.Email,
		Primary:  utils.IfNotNilBool(input.Primary),
		Source:   neo4jentity.DataSourceOpenline,
	}
	return &emailEntity
}

func MapEntityToEmail(entity *neo4jentity.EmailEntity) *model.Email {
	if entity == nil {
		return nil
	}
	return &model.Email{
		ID:        entity.Id,
		Email:     utils.StringPtrFirstNonEmptyNillable(entity.Email, entity.RawEmail),
		RawEmail:  utils.StringPtrNillable(entity.RawEmail),
		Work:      entity.Work,
		Primary:   entity.Primary,
		Source:    MapDataSourceToModel(entity.Source),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		EmailValidationDetails: &model.EmailValidationDetails{
			Verified:          entity.EmailInternalFields.ValidatedAt != nil,
			VerifyingCheckAll: false,
			IsValidSyntax:     entity.IsValidSyntax,
			IsCatchAll:        entity.IsCatchAll,
			IsRoleAccount:     entity.IsRoleAccount,
			IsSystemGenerated: entity.IsSystemGenerated,
			IsRisky:           entity.IsRisky,
			IsFirewalled:      entity.IsFirewalled,
			Provider:          entity.Provider,
			Firewall:          entity.Firewall,
			IsMailboxFull:     entity.IsMailboxFull,
			IsFreeAccount:     entity.IsFreeAccount,
			SMTPSuccess:       entity.SmtpSuccess,
			Deliverable:       enummapper.MapDeliverableToModelPtr(entity.Deliverable),
			IsPrimaryDomain:   entity.IsPrimaryDomain,
			PrimaryDomain:     entity.PrimaryDomain,
			AlternateEmail:    entity.AlternateEmail,
		},
	}
}

func MapEmailInputToLocalEntity(input *model.EmailInput) *neo4jentity.EmailEntity {
	if input == nil {
		return nil
	}
	emailEntity := neo4jentity.EmailEntity{
		RawEmail: input.Email,
		Primary:  utils.IfNotNilBool(input.Primary),
		Source:   neo4jentity.DataSourceOpenline,
	}
	return &emailEntity
}

func MapEntitiesToEmails(entities *neo4jentity.EmailEntities) []*model.Email {
	if entities == nil {
		return nil
	}
	emails := make([]*model.Email, 0, len(*entities))
	for _, entity := range *entities {
		emails = append(emails, MapEntityToEmail(&entity))
	}
	return emails
}
