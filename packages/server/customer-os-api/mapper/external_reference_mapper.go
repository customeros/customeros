package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapExternalSystemReferenceInputToRelationship(input *model.ExternalSystemReferenceInput) *neo4jentity.ExternalSystemEntity {
	if input == nil {
		return nil
	}
	relationship := new(neo4jentity.ExternalSystemEntity)
	relationship.Relationship.ExternalId = input.ExternalID
	relationship.Relationship.ExternalUrl = input.ExternalURL
	relationship.Relationship.ExternalSource = input.ExternalSource
	relationship.ExternalSystemId = mapper.MapExternalSystemTypeFromModel(input.Type)

	if input.SyncDate == nil {
		relationship.Relationship.SyncDate = utils.ToPtr(utils.Now())
	} else {
		relationship.Relationship.SyncDate = input.SyncDate
	}
	return relationship
}

func MapEntityToExternalSystem(entity *neo4jentity.ExternalSystemEntity) *model.ExternalSystem {
	return &model.ExternalSystem{
		Type:           mapper.MapExternalSystemTypeToModel(entity.ExternalSystemId),
		SyncDate:       entity.Relationship.SyncDate,
		ExternalID:     utils.StringPtrNillable(entity.Relationship.ExternalId),
		ExternalURL:    entity.Relationship.ExternalUrl,
		ExternalSource: entity.Relationship.ExternalSource,
	}
}

func MapEntitiesToExternalSystems(entities *neo4jentity.ExternalSystemEntities) []*model.ExternalSystem {
	var output []*model.ExternalSystem
	for _, v := range *entities {
		output = append(output, MapEntityToExternalSystem(&v))
	}
	return output
}
