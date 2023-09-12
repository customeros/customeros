package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var externalSystemTypeByModel = map[model.ExternalSystemType]entity.ExternalSystemId{
	model.ExternalSystemTypeHubspot:        entity.Hubspot,
	model.ExternalSystemTypeZendeskSupport: entity.ZendeskSupport,
	model.ExternalSystemTypeCalcom:         entity.CalCom,
	model.ExternalSystemTypePipedrive:      entity.Pipedrive,
	model.ExternalSystemTypeSLACk:          entity.Slack,
	model.ExternalSystemTypeIntercom:       entity.Intercom,
	model.ExternalSystemTypeSalesforce:     entity.Salesforce,
}

var externalSystemTypeByValue = utils.ReverseMap(externalSystemTypeByModel)

func MapExternalSystemTypeFromModel(input model.ExternalSystemType) entity.ExternalSystemId {
	return externalSystemTypeByModel[input]
}

func MapExternalSystemTypeToModel(input entity.ExternalSystemId) model.ExternalSystemType {
	return externalSystemTypeByValue[input]
}
