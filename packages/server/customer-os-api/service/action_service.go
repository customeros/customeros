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
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

type ActionService interface {
	GetActionsForNodes(ctx context.Context, entityType neo4jenum.EntityType, ids []string) (*entity.ActionEntities, error)

	mapDbNodeToActionEntity(dbNode dbtype.Node) *entity.ActionEntity
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

func (s *actionService) GetActionsForNodes(ctx context.Context, entityType neo4jenum.EntityType, ids []string) (*entity.ActionEntities, error) {
	records, err := s.repositories.Neo4jRepositories.ActionReadRepository.GetFor(ctx, common.GetTenantFromContext(ctx), entityType, ids)
	if err != nil {
		return nil, err
	}

	var data entity.ActionEntities
	for _, v := range records {
		action := s.mapDbNodeToActionEntity(*v.Node)
		action.DataloaderKey = v.LinkedNodeId
		data = append(data, *action)
	}

	return &data, nil
}

func (s *actionService) mapDbNodeToActionEntity(dbNode dbtype.Node) *entity.ActionEntity {
	props := utils.GetPropsFromNode(dbNode)
	action := entity.ActionEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Type:      neo4jenum.GetActionType(utils.GetStringPropOrEmpty(props, "type")),
		Content:   utils.GetStringPropOrEmpty(props, "content"),
		Metadata:  utils.GetStringPropOrEmpty(props, "metadata"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		AppSource: utils.GetStringPropOrEmpty(props, "appSource"),
		Source:    neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
	}
	return &action
}
