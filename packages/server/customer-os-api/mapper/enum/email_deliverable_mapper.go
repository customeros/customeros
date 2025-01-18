package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/verify"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

var deliverableByModel = map[model.EmailDeliverable]string{
	model.EmailDeliverableDeliverable:   string(verify.EmailDeliverableStatusDeliverable),
	model.EmailDeliverableUndeliverable: string(verify.EmailDeliverableStatusUndeliverable),
	model.EmailDeliverableUnknown:       string(verify.EmailDeliverableStatusUnknown),
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
