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
	return &emailEntity
}
