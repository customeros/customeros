package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var externalSystemTypeByModel = map[model.ExternalSystemType]neo4jenum.ExternalSystemId{
	model.ExternalSystemTypeHubspot:        neo4jenum.Hubspot,
	model.ExternalSystemTypeZendeskSupport: neo4jenum.ZendeskSupport,
	model.ExternalSystemTypeCalcom:         neo4jenum.CalCom,
	model.ExternalSystemTypePipedrive:      neo4jenum.Pipedrive,
	model.ExternalSystemTypeSLACk:          neo4jenum.Slack,
	model.ExternalSystemTypeIntercom:       neo4jenum.Intercom,
	model.ExternalSystemTypeSalesforce:     neo4jenum.Salesforce,
	model.ExternalSystemTypeStripe:         neo4jenum.Stripe,
	model.ExternalSystemTypeMixpanel:       neo4jenum.Mixpanel,
	model.ExternalSystemTypeClose:          neo4jenum.Close,
	model.ExternalSystemTypeOutlook:        neo4jenum.Outlook,
	model.ExternalSystemTypeUnthread:       neo4jenum.Unthread,
	model.ExternalSystemTypeAttio:          neo4jenum.Attio,
	model.ExternalSystemTypeWeconnect:      neo4jenum.WeConnect,
}

var externalSystemTypeByValue = utils.ReverseMap(externalSystemTypeByModel)

func MapExternalSystemTypeFromModel(input model.ExternalSystemType) neo4jenum.ExternalSystemId {
	return externalSystemTypeByModel[input]
}

func MapExternalSystemTypeToModel(input neo4jenum.ExternalSystemId) model.ExternalSystemType {
	return externalSystemTypeByValue[input]
}
