package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"time"
)

func MapExternalSystemReferenceInputToRelationship(input *model.ExternalSystemReferenceInput) *entity.ExternalReferenceRelationship {
	if input == nil {
		return nil
	}
	relationship := new(entity.ExternalReferenceRelationship)
	relationship.Id = input.ID
	relationship.ExternalSystemId = MapExternalSystemTypeFromModel(input.Type)
	if input.SyncDate == nil {
		relationship.SyncDate = time.Now().UTC()
	} else {
		relationship.SyncDate = *input.SyncDate
	}
	return relationship
}
