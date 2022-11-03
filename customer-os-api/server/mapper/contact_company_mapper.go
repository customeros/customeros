package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapContactCompanyInputToEntity(input *model.ContactCompanyInput) *entity.ContactCompanyEntity {
	if input == nil {
		return nil
	}
	contactCompanyEntity := entity.ContactCompanyEntity{
		Company: input.CompanyName,
	}
	if input.JobTitle != nil {
		contactCompanyEntity.JobTitle = *input.JobTitle
	}
	return &contactCompanyEntity
}
