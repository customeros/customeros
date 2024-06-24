package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SocialService interface {
	Update(ctx context.Context, entity neo4jentity.SocialEntity) (*neo4jentity.SocialEntity, error)
	GetAllForEntities(ctx context.Context, linkedEntityType entity.EntityType, linkedEntityIds []string) (*neo4jentity.SocialEntities, error)
	Remove(ctx context.Context, socialId string) error
}

type socialService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewSocialService(log logger.Logger, repositories *repository.Repositories) SocialService {
	return &socialService{
		log:          log,
		repositories: repositories,
	}
}

func (s *socialService) GetAllForEntities(ctx context.Context, linkedEntityType entity.EntityType, linkedEntityIds []string) (*neo4jentity.SocialEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.GetAllForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("linkedEntityType", string(linkedEntityType)), log.Object("linkedEntityIds", linkedEntityIds))

	socials, err := s.repositories.SocialRepository.GetAllForEntities(ctx, common.GetTenantFromContext(ctx), linkedEntityType, linkedEntityIds)
	if err != nil {
		return nil, err
	}
	socialEntities := make(neo4jentity.SocialEntities, 0, len(socials))
	for _, v := range socials {
		socialEntity := s.mapDbNodeToSocialEntity(*v.Node)
		socialEntity.DataloaderKey = v.LinkedNodeId
		socialEntities = append(socialEntities, *socialEntity)
	}
	return &socialEntities, nil
}

func (s *socialService) Update(ctx context.Context, socialEntity neo4jentity.SocialEntity) (*neo4jentity.SocialEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	updatedLocationNode, err := s.repositories.SocialRepository.Update(ctx, common.GetTenantFromContext(ctx), socialEntity)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToSocialEntity(*updatedLocationNode), nil
}

func (s *socialService) Remove(ctx context.Context, socialId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.Remove")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, socialId)

	return s.repositories.Neo4jRepositories.SocialWriteRepository.PermanentlyDelete(ctx, common.GetTenantFromContext(ctx), socialId)
}

func (s *socialService) mapDbNodeToSocialEntity(node dbtype.Node) *neo4jentity.SocialEntity {
	props := utils.GetPropsFromNode(node)
	return &neo4jentity.SocialEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Url:           utils.GetStringPropOrEmpty(props, "url"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
}
