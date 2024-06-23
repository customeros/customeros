package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
)

type ActionService interface {
	GetActionsForNodes(ctx context.Context, entityType neo4jenum.EntityType, ids []string) (*neo4jentity.ActionEntities, error)
}

type actionService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewActionService(log logger.Logger, repository *repository.Repositories) ActionService {
	return &actionService{
		log:          log,
		repositories: repository,
	}
}

func (s *actionService) GetActionsForNodes(ctx context.Context, entityType neo4jenum.EntityType, ids []string) (*neo4jentity.ActionEntities, error) {
	records, err := s.repositories.Neo4jRepositories.ActionReadRepository.GetFor(ctx, common.GetTenantFromContext(ctx), entityType, ids)
	if err != nil {
		return nil, err
	}

	var data neo4jentity.ActionEntities
	for _, v := range records {
		action := neo4jmapper.MapDbNodeToActionEntity(v.Node)
		action.DataloaderKey = v.LinkedNodeId
		data = append(data, *action)
	}

	return &data, nil
}
