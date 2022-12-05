package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"golang.org/x/exp/slices"
	"time"
)

type ActionsService interface {
	GetContactActions(context.Context, string, time.Time, time.Time, []model.ActionType) (*entity.ActionEntities, error)
}

type actionsService struct {
	repositories *repository.Repositories
}

func NewActionsService(repositories *repository.Repositories) ActionsService {
	return &actionsService{
		repositories: repositories,
	}
}

func (s *actionsService) getNeo4jDriver() neo4j.Driver {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *actionsService) GetContactActions(ctx context.Context, contactId string, from time.Time, to time.Time, types []model.ActionType) (*entity.ActionEntities, error) {
	session := utils.NewNeo4jReadSession(s.getNeo4jDriver())
	defer session.Close()

	var nodeLabels = []string{}
	for _, v := range types {
		nodeLabels = append(nodeLabels, entity.NodeLabelsByActionType[v.String()])
	}

	dbNodes, err := s.repositories.ActionRepository.GetContactActions(session, common.GetContext(ctx).Tenant, contactId, from, to, nodeLabels)
	if err != nil {
		return nil, err
	}

	actions := entity.ActionEntities{}
	for _, v := range dbNodes {
		if slices.Contains(v.Labels, entity.LabelName_PageViewAction) {
			actions = append(actions, s.mapDbNodeToPageViewAction(v))
		}
	}

	return &actions, nil
}

func (s *actionsService) mapDbNodeToPageViewAction(node *dbtype.Node) *entity.PageViewActionEntity {
	props := utils.GetPropsFromNode(*node)
	pageViewAction := entity.PageViewActionEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		Application:    utils.GetStringPropOrEmpty(props, "application"),
		TrackerName:    utils.GetStringPropOrEmpty(props, "trackerName"),
		SessionId:      utils.GetStringPropOrEmpty(props, "sessionId"),
		PageUrl:        utils.GetStringPropOrEmpty(props, "pageUrl"),
		PageTitle:      utils.GetStringPropOrEmpty(props, "pageTitle"),
		OrderInSession: utils.GetIntPropOrZero(props, "orderInSession"),
		EngagedTime:    utils.GetIntPropOrZero(props, "engagedTime"),
		StartedAt:      utils.GetTimePropOrNow(props, "startedAt"),
		EndedAt:        utils.GetTimePropOrNow(props, "endedAt"),
	}
	return &pageViewAction
}
