package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type ActionItemService interface {
	GetActionItemsForNodes(ctx context.Context, linkedWith repository.LinkedWith, ids []string) (*entity.ActionItemEntities, error)

	MapDbNodeToActionItemEntity(node dbtype.Node) *entity.ActionItemEntity
}

type actionItemService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewActionItemService(log logger.Logger, repositories *repository.Repositories) ActionItemService {
	return &actionItemService{
		log:          log,
		repositories: repositories,
	}
}

func (s *actionItemService) GetActionItemsForNodes(ctx context.Context, linkedWith repository.LinkedWith, ids []string) (*entity.ActionItemEntities, error) {
	records, err := s.repositories.ActionItemRepository.GetFor(ctx, common.GetTenantFromContext(ctx), linkedWith, ids)
	if err != nil {
		return nil, err
	}

	converted := s.convertDbNodesToActionItems(records)

	return &converted, nil
}

func (s *actionItemService) convertDbNodesToActionItems(records []*utils.DbNodeAndId) entity.ActionItemEntities {
	entities := entity.ActionItemEntities{}
	for _, v := range records {
		entity := s.MapDbNodeToActionItemEntity(*v.Node)
		entity.DataloaderKey = v.LinkedNodeId
		entities = append(entities, *entity)

	}
	return entities
}

func (s *actionItemService) MapDbNodeToActionItemEntity(node dbtype.Node) *entity.ActionItemEntity {
	props := utils.GetPropsFromNode(node)
	createdAt := utils.GetTimePropOrEpochStart(props, "createdAt")
	entity := entity.ActionItemEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     &createdAt,
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4jentity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &entity
}
