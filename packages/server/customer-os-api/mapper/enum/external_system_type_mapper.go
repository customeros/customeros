package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

var externalSystemTypeByModel = map[model.ExternalSystemType]enum.Source{
	model.ExternalSystemTypeAttio:          enum.SourceAttio,
	model.ExternalSystemTypeCalcom:         enum.SourceCalCom,
	model.ExternalSystemTypeClose:          enum.SourceClose,
	model.ExternalSystemTypeHubspot:        enum.SourceHubspot,
	model.ExternalSystemTypeZendeskSupport: enum.SourceZendeskSupport,
	model.ExternalSystemTypePipedrive:      enum.SourcePipedrive,
	model.ExternalSystemTypeSLACk:          enum.SourceSlack,
	model.ExternalSystemTypeIntercom:       enum.SourceIntercom,
	model.ExternalSystemTypeSalesforce:     enum.SourceSalesforce,
	model.ExternalSystemTypeStripe:         enum.SourceStripe,
	model.ExternalSystemTypeMixpanel:       enum.SourceMixpanel,
	model.ExternalSystemTypeOutlook:        enum.SourceOutlook,
	model.ExternalSystemTypeUnthread:       enum.SourceUnthread,
	model.ExternalSystemTypeZendeskSell:    enum.SourceZendeskSell,
}

var externalSystemTypeByValue = utils.ReverseMap(externalSystemTypeByModel)

func MapExternalSystemTypeFromModel(input model.ExternalSystemType) enum.Source {
	return externalSystemTypeByModel[input]
}

func MapExternalSystemTypeToModel(input enum.Source) model.ExternalSystemType {
	return externalSystemTypeByValue[input]
}
