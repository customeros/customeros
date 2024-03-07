package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/pkg/errors"
	"strings"
)

func GetTagId(ctx context.Context, services *service.Services, tagId, tagName *string) string {
	outputTagId := ""
	if tagId != nil && *tagId != "" {
		tagEntity, _ := services.TagService.GetById(ctx, *tagId)
		if tagEntity != nil {
			outputTagId = tagEntity.Id
		}
	}
	if outputTagId == "" && tagName != nil && strings.TrimSpace(*tagName) != "" {
		tagEntity, _ := services.TagService.GetByNameOptional(ctx, strings.TrimSpace(*tagName))
		if tagEntity != nil {
			outputTagId = tagEntity.Id
		}
	}
	return outputTagId
}

func CreateTag(ctx context.Context, services *service.Services, tagName *string) (*neo4jentity.TagEntity, error) {
	if tagName == nil || strings.TrimSpace(*tagName) == "" {
		return nil, errors.New("tag name is empty")
	}
	return services.TagService.Merge(ctx, &neo4jentity.TagEntity{
		Name:          strings.TrimSpace(*tagName),
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     constants.AppSourceCustomerOsApi,
	})
}
