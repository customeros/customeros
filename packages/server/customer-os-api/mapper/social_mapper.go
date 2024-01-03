package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"strings"
)

func MapSocialInputToEntity(input *model.SocialInput) *entity.SocialEntity {
	if input == nil {
		return &entity.SocialEntity{}
	}
	return &entity.SocialEntity{
		SourceFields: entity.SourceFields{
			Source:        neo4jentity.DataSourceOpenline,
			SourceOfTruth: neo4jentity.DataSourceOpenline,
			AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
		},
		PlatformName: utils.IfNotNilString(input.PlatformName, func() string { return strings.ToLower(*input.PlatformName) }),
		Url:          input.URL,
	}
}

func MapSocialUpdateInputToEntity(input *model.SocialUpdateInput) *entity.SocialEntity {
	if input == nil {
		return &entity.SocialEntity{}
	}
	return &entity.SocialEntity{
		Id: input.ID,
		SourceFields: entity.SourceFields{
			SourceOfTruth: neo4jentity.DataSourceOpenline,
		},
		PlatformName: utils.IfNotNilString(input.PlatformName, func() string { return strings.ToLower(*input.PlatformName) }),
		Url:          utils.IfNotNilString(input.URL),
	}
}

func MapEntityToSocial(entity *entity.SocialEntity) *model.Social {
	return &model.Social{
		ID:            entity.Id,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		PlatformName:  utils.StringPtr(entity.PlatformName),
		URL:           entity.Url,
		Source:        MapDataSourceToModel(entity.SourceFields.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceFields.SourceOfTruth),
		AppSource:     entity.SourceFields.AppSource,
	}
}

func MapEntitiesToSocials(entities *entity.SocialEntities) []*model.Social {
	var socials []*model.Social
	for _, socialEntity := range *entities {
		socials = append(socials, MapEntityToSocial(&socialEntity))
	}
	return socials
}
