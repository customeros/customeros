package mapper

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapDbNodeToMasterPlanEntity(dbNode *dbtype.Node) *neo4jentity.MasterPlanEntity {
	if dbNode == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*dbNode)
	entity := neo4jentity.MasterPlanEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &entity
}

func MapDbNodeToMasterPlanMilestoneEntity(dbNode *dbtype.Node) *neo4jentity.MasterPlanMilestoneEntity {
	if dbNode == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*dbNode)
	entity := neo4jentity.MasterPlanMilestoneEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Order:         utils.GetInt64PropOrZero(props, "order"),
		DurationHours: utils.GetInt64PropOrZero(props, "durationHours"),
		Optional:      utils.GetBoolPropOrFalse(props, "optional"),
		Items:         utils.GetListStringPropOrEmpty(props, "items"),
	}
	return &entity
}
