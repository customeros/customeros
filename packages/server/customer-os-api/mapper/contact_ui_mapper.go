package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToContactUIDetails(entity *neo4jentity.ContactEntity, output *model.ContactUIDetails) {
	if entity == nil || output == nil {
		return
	}

	output.ID = entity.Id
	output.CreatedAt = entity.CreatedAt
	output.UpdatedAt = entity.UpdatedAt
	output.Hide = entity.Hide
	output.FirstName = entity.FirstName
	output.LastName = entity.LastName
	output.Name = entity.Name
	output.Prefix = entity.Prefix
	output.Description = entity.Description
	output.Timezone = entity.Timezone
	output.ProfilePhotoURL = entity.ProfilePhotoUrl
}
