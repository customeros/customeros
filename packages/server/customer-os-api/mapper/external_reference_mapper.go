package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapExternalSystemReferenceInputToRelationship(input *model.ExternalSystemReferenceInput) *entity.ExternalSystemEntity {
	if input == nil {
		return nil
	}
	relationship := new(entity.ExternalSystemEntity)
	relationship.Relationship.ExternalId = input.ExternalID
	relationship.Relationship.ExternalUrl = input.ExternalURL
	relationship.Relationship.ExternalSource = input.ExternalSource
	relationship.ExternalSystemId = MapExternalSystemTypeFromModel(input.Type)

	if input.SyncDate == nil {
		relationship.Relationship.SyncDate = utils.ToPtr(utils.Now())
	} else {
		relationship.Relationship.SyncDate = input.SyncDate
	}
	return relationship
}

func MapEntityToExternalSystem(entity *entity.ExternalSystemEntity) *model.ExternalSystem {
	return &model.ExternalSystem{
		Type:           MapExternalSystemTypeToModel(entity.ExternalSystemId),
		SyncDate:       entity.Relationship.SyncDate,
		ExternalID:     utils.StringPtrNillable(entity.Relationship.ExternalId),
		ExternalURL:    entity.Relationship.ExternalUrl,
		ExternalSource: entity.Relationship.ExternalSource,
	}
}

func MapEntitiesToExternalSystems(entities *entity.ExternalSystemEntities) []*model.ExternalSystem {
	var output []*model.ExternalSystem
	for _, v := range *entities {
		output = append(output, MapEntityToExternalSystem(&v))
	}
	return output
}
