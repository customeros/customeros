package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type SocialService interface {
	CreateSocialForEntity(ctx context.Context, linkedEntityType entity.EntityType, linkedEntityId string, socialEntity entity.SocialEntity) (*entity.SocialEntity, error)
	Update(ctx context.Context, entity entity.SocialEntity) (*entity.SocialEntity, error)
	GetAllForEntities(ctx context.Context, linkedEntityType entity.EntityType, linkedEntityIds []string) (*entity.SocialEntities, error)
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

func (s *socialService) GetAllForEntities(ctx context.Context, linkedEntityType entity.EntityType, linkedEntityIds []string) (*entity.SocialEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.GetAllForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("linkedEntityType", string(linkedEntityType)), log.Object("linkedEntityIds", linkedEntityIds))

	socials, err := s.repositories.SocialRepository.GetAllForEntities(ctx, common.GetTenantFromContext(ctx), linkedEntityType, linkedEntityIds)
	if err != nil {
		return nil, err
	}
	socialEntities := make(entity.SocialEntities, 0, len(socials))
	for _, v := range socials {
		socialEntity := s.mapDbNodeToSocialEntity(*v.Node)
		socialEntity.DataloaderKey = v.LinkedNodeId
		socialEntities = append(socialEntities, *socialEntity)
	}
	return &socialEntities, nil
}

func (s *socialService) CreateSocialForEntity(ctx context.Context, linkedEntityType entity.EntityType, linkedEntityId string, socialEntity entity.SocialEntity) (*entity.SocialEntity, error) {
	if linkedEntityType != entity.CONTACT && linkedEntityType != entity.ORGANIZATION {
		return nil, errors.ErrInvalidEntityType
	}
	socialNode, err := s.repositories.SocialRepository.CreateSocialForEntity(ctx, common.GetTenantFromContext(ctx), linkedEntityType, linkedEntityId, socialEntity)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToSocialEntity(*socialNode), nil
}

func (s *socialService) Update(ctx context.Context, socialEntity entity.SocialEntity) (*entity.SocialEntity, error) {
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("socialId", socialId))

	return s.repositories.SocialRepository.Remove(ctx, socialId)
}

func (s *socialService) mapDbNodeToSocialEntity(node dbtype.Node) *entity.SocialEntity {
	props := utils.GetPropsFromNode(node)
	return &entity.SocialEntity{
		Id:           utils.GetStringPropOrEmpty(props, "id"),
		PlatformName: utils.GetStringPropOrEmpty(props, "platformName"),
		Url:          utils.GetStringPropOrEmpty(props, "url"),
		CreatedAt:    utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:    utils.GetTimePropOrEpochStart(props, "updatedAt"),
		SourceFields: entity.SourceFields{
			Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
			SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
			AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		},
	}
}
