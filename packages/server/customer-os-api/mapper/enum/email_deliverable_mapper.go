package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
)

var deliverableByModel = map[model.EmailDeliverable]string{
	model.EmailDeliverableDeliverable:   string(validationmodel.EmailDeliverableStatusDeliverable),
	model.EmailDeliverableUndeliverable: string(validationmodel.EmailDeliverableStatusUndeliverable),
	model.EmailDeliverableUnknown:       string(validationmodel.EmailDeliverableStatusUnknown),
}

var deliverableByValue = utils.ReverseMap(deliverableByModel)

func MapDeliverableToModelPtr(input *string) *model.EmailDeliverable {
	if input == nil {
		return nil
	}
	if v, exists := deliverableByValue[*input]; exists {
		return &v
	} else {
		return utils.ToPtr(model.EmailDeliverableUnknown)
	}
}
