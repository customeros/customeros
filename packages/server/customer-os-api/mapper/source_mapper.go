package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var sourceByModel = map[model.DataSource]entity.DataSource{
	model.DataSourceNa:             entity.DataSourceNA,
	model.DataSourceOpenline:       entity.DataSourceOpenline,
	model.DataSourceHubspot:        entity.DataSourceHubspot,
	model.DataSourceZendeskSupport: entity.DataSourceZendeskSupport,
	model.DataSourcePipedrive:      entity.DataSourcePipedrive,
	model.DataSourceNotion:         entity.DataSourceNotion,
	model.DataSourceSLACk:          entity.DataSourceSlack,
	model.DataSourceWebscrape:      entity.DataSourceWebscrape,
}

var sourceByValue = utils.ReverseMap(sourceByModel)

func MapDataSourceFromModel(input model.DataSource) entity.DataSource {
	return sourceByModel[input]
}

func MapDataSourceToModel(input entity.DataSource) model.DataSource {
	return sourceByValue[input]
}
