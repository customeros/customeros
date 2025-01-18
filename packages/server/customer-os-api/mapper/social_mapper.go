package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapSocialUpdateInputToEntity(input *model.SocialUpdateInput) *neo4jentity.SocialEntity {
	if input == nil {
		return &neo4jentity.SocialEntity{}
	}
	return &neo4jentity.SocialEntity{
		Id:  input.ID,
		Url: utils.IfNotNilString(input.URL),
	}
}

func MapEntityToSocial(entity *neo4jentity.SocialEntity) *model.Social {
	return &model.Social{
		Metadata: &model.Metadata{
			ID:          entity.Id,
			Created:     entity.CreatedAt,
			LastUpdated: entity.UpdatedAt,
			Source:      MapDataSourceToModel(entity.Source),
			AppSource:   entity.AppSource,
		},
		ID:             entity.Id,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		URL:            entity.Url,
		Alias:          entity.Alias,
		FollowersCount: entity.FollowersCount,
		ExternalID:     entity.ExternalId,
		Source:         MapDataSourceToModel(entity.Source),
		AppSource:      entity.AppSource,
	}
}

func MapEntitiesToSocials(entities *neo4jentity.SocialEntities) []*model.Social {
	var socials []*model.Social
	for _, socialEntity := range *entities {
		socials = append(socials, MapEntityToSocial(&socialEntity))
	}
	return socials
}
