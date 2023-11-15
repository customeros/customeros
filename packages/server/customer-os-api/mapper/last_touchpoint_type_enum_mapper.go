package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var lastTouchpointTypeModel = map[model.LastTouchpointType]string{
	model.LastTouchpointTypePageView:                  string(entity.LastTouchpointTypePageView),
	model.LastTouchpointTypeInteractionSession:        string(entity.LastTouchpointTypeInteractionSession),
	model.LastTouchpointTypeNote:                      string(entity.LastTouchpointTypeNote),
	model.LastTouchpointTypeInteractionEventEmailSent: string(entity.LastTouchpointTypeInteractionEventEmailSent),
	model.LastTouchpointTypeInteractionEventPhoneCall: string(entity.LastTouchpointTypeInteractionEventPhoneCall),
	model.LastTouchpointTypeInteractionEventChat:      string(entity.LastTouchpointTypeInteractionEventChat),
	model.LastTouchpointTypeMeeting:                   string(entity.LastTouchpointTypeMeeting),
	model.LastTouchpointTypeAnalysis:                  string(entity.LastTouchpointTypeAnalysis),
	model.LastTouchpointTypeActionCreated:             string(entity.LastTouchpointTypeActionCreated),
	model.LastTouchpointTypeAction:                    string(entity.LastTouchpointTypeAction),
	model.LastTouchpointTypeLogEntry:                  string(entity.LastTouchpointTypeLogEntry),
	model.LastTouchpointTypeIssueCreated:              string(entity.LastTouchpointTypeIssueCreated),
	model.LastTouchpointTypeIssueUpdated:              string(entity.LastTouchpointTypeIssueUpdated),
}

var lastTouchpointTypeByValue = utils.ReverseMap(lastTouchpointTypeModel)

func MapLastTouchpointTypeFromModel(input *model.LastTouchpointType) string {
	if input == nil {
		return ""
	}
	if v, exists := lastTouchpointTypeModel[*input]; exists {
		return v
	} else {
		return ""
	}
}

func MapLastTouchpointTypeFromString(input *string) string {
	if input == nil {
		return ""
	}
	if v, exists := lastTouchpointTypeModel[model.LastTouchpointType(*input)]; exists {
		return v
	} else {
		return ""
	}
}

func MapLastTouchpointTypeToModel(input *string) *model.LastTouchpointType {
	if input == nil {
		return nil
	}
	if v, exists := lastTouchpointTypeByValue[*input]; exists {
		return &v
	} else {
		return nil
	}
}
