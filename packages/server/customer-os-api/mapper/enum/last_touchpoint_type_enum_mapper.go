package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var lastTouchpointTypeModel = map[model.LastTouchpointType]string{
	model.LastTouchpointTypePageView:                      string(neo4jenum.TouchpointTypePageView),
	model.LastTouchpointTypeInteractionSession:            string(neo4jenum.TouchpointTypeInteractionSession),
	model.LastTouchpointTypeNote:                          string(neo4jenum.TouchpointTypeNote),
	model.LastTouchpointTypeInteractionEventEmailSent:     string(neo4jenum.TouchpointTypeInteractionEventEmailSent),
	model.LastTouchpointTypeInteractionEventEmailReceived: string(neo4jenum.TouchpointTypeInteractionEventEmailReceived),
	model.LastTouchpointTypeInteractionEventPhoneCall:     string(neo4jenum.TouchpointTypeInteractionEventPhoneCall),
	model.LastTouchpointTypeInteractionEventChat:          string(neo4jenum.TouchpointTypeInteractionEventChat),
	model.LastTouchpointTypeMeeting:                       string(neo4jenum.TouchpointTypeMeeting),
	model.LastTouchpointTypeActionCreated:                 string(neo4jenum.TouchpointTypeActionCreated),
	model.LastTouchpointTypeAction:                        string(neo4jenum.TouchpointTypeAction),
	model.LastTouchpointTypeLogEntry:                      string(neo4jenum.TouchpointTypeLogEntry),
	model.LastTouchpointTypeIssueCreated:                  string(neo4jenum.TouchpointTypeIssueCreated),
	model.LastTouchpointTypeIssueUpdated:                  string(neo4jenum.TouchpointTypeIssueUpdated),
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
