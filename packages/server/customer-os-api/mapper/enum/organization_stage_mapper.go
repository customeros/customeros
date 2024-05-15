package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var stageByModel = map[model.OrganizationStage]neo4jenum.OrganizationStage{
	model.OrganizationStageLead:        neo4jenum.Lead,
	model.OrganizationStageTarget:      neo4jenum.Target,
	model.OrganizationStageInterested:  neo4jenum.Interested,
	model.OrganizationStageEngaged:     neo4jenum.Engaged,
	model.OrganizationStageClosedLost:  neo4jenum.ClosedLost,
	model.OrganizationStageClosedWon:   neo4jenum.ClosedWon,
	model.OrganizationStageUnqualified: neo4jenum.Unqualified,
	model.OrganizationStageNurture:     neo4jenum.Nurture,
}

var stageByValue = utils.ReverseMap(stageByModel)

func MapStageFromModel(input model.OrganizationStage) neo4jenum.OrganizationStage {
	return stageByModel[input]
}

func MapStageToModel(input neo4jenum.OrganizationStage) model.OrganizationStage {
	return stageByValue[input]
}
