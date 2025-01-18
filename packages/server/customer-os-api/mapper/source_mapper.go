package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var sourceByModel = map[model.DataSource]neo4jentity.DataSource{
	model.DataSourceNa:             neo4jentity.DataSourceNA,
	model.DataSourceOpenline:       neo4jentity.DataSourceOpenline,
	model.DataSourceHubspot:        neo4jentity.DataSourceHubspot,
	model.DataSourceZendeskSupport: neo4jentity.DataSourceZendeskSupport,
	model.DataSourcePipedrive:      neo4jentity.DataSourcePipedrive,
	model.DataSourceSLACk:          neo4jentity.DataSourceSlack,
	model.DataSourceWebscrape:      neo4jentity.DataSourceWebscrape,
	model.DataSourceIntercom:       neo4jentity.DataSourceIntercom,
	model.DataSourceSalesforce:     neo4jentity.DataSourceSalesforce,
	model.DataSourceStripe:         neo4jentity.DataSourceStripe,
	model.DataSourceMixpanel:       neo4jentity.DataSourceMixpanel,
	model.DataSourceClose:          neo4jentity.DataSourceClose,
	model.DataSourceOutlook:        neo4jentity.DataSourceOutlook,
	model.DataSourceUnthread:       neo4jentity.DataSourceUnthread,
	model.DataSourceAttio:          neo4jentity.DataSourceAttio,
	model.DataSourceZendeskSell:    neo4jentity.DataSourceZendeskSell,
	model.DataSourceMailstack:      neo4jentity.DataSourceMailstack,
	model.DataSourceFathom:         neo4jentity.DataSourceFathom,
	model.DataSourceGrain:          neo4jentity.DataSourceGrain,
}

var sourceByValue = utils.ReverseMap(sourceByModel)

func MapDataSourceFromModel(input model.DataSource) neo4jentity.DataSource {
	return sourceByModel[input]
}

func MapDataSourceToModel(input neo4jentity.DataSource) model.DataSource {
	return sourceByValue[input]
}
