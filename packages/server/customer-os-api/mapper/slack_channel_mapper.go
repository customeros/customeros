package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
)

func MapEntityToSlackChannel(entity *postgresEntity.SlackChannel) *model.SlackChannel {
	if entity == nil {
		return nil
	}
	output := model.SlackChannel{
		Metadata: &model.Metadata{
			ID:          entity.ID.String(),
			Created:     entity.CreatedAt,
			LastUpdated: entity.UpdatedAt,
			Source:      model.DataSource(entity.Source),
		},
		ChannelID: entity.ChannelId,
	}
	return &output
}

func MapEntitiesToSlackChannels(entities []*postgresEntity.SlackChannel) []*model.SlackChannel {
	var output []*model.SlackChannel
	for _, v := range entities {
		output = append(output, MapEntityToSlackChannel(v))
	}
	return output
}
