package mapper

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

const (
	typeHubspot        = "hubspot"
	typeZendeskSupport = "zendesk_support"
)

var externalSystemTypeByModel = map[model.ExternalSystemType]string{
	model.ExternalSystemTypeHubspot:        typeHubspot,
	model.ExternalSystemTypeZendeskSupport: typeZendeskSupport,
}

var externalSystemTypeByValue = map[string]model.ExternalSystemType{
	typeHubspot:        model.ExternalSystemTypeHubspot,
	typeZendeskSupport: model.ExternalSystemTypeZendeskSupport,
}

func MapExternalSystemTypeFromModel(input model.ExternalSystemType) string {
	return externalSystemTypeByModel[input]
}

func MapExternalSystemTypeToModel(input string) model.ExternalSystemType {
	return externalSystemTypeByValue[input]
}
