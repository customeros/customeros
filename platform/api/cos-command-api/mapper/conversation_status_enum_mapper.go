package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

const (
	statusActive = "ACTIVE"
	statusClosed = "CLOSED"
)

var statusByModel = map[model.ConversationStatus]string{
	model.ConversationStatusActive: statusActive,
	model.ConversationStatusClosed: statusClosed,
}

var statusByValue = utils.ReverseMap(statusByModel)

func MapConversationStatusFromModel(input model.ConversationStatus) string {
	return statusByModel[input]
}

func MapConversationStatusToModel(input string) model.ConversationStatus {
	return statusByValue[input]
}
