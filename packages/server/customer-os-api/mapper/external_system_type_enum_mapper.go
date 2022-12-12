package mapper

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"

const (
	typeHubspot = "hubspot"
	typeZendesk = "zendesk"
)

var externalSystemTypeByModel = map[model.ExternalSystemType]string{
	model.ExternalSystemTypeHubspot: typeHubspot,
	model.ExternalSystemTypeZendesk: typeZendesk,
}

var externalSystemTypeByValue = map[string]model.ExternalSystemType{
	typeHubspot: model.ExternalSystemTypeHubspot,
	typeZendesk: model.ExternalSystemTypeZendesk,
}

func MapExternalSystemTypeFromModel(input model.ExternalSystemType) string {
	return externalSystemTypeByModel[input]
}

func MapExternalSystemTypeToModel(input string) model.ExternalSystemType {
	return externalSystemTypeByValue[input]
}
