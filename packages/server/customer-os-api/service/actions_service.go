package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/exp/slices"
	"time"
)

type ActionsService interface {
	GetContactActions(context.Context, string, time.Time, time.Time, []model.ActionType) (*entity.ActionEntities, error)
}

type actionsService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewActionsService(repositories *repository.Repositories, services *Services) ActionsService {
	return &actionsService{
		repositories: repositories,
		services:     services,
	}
}

func (s *actionsService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *actionsService) GetContactActions(ctx context.Context, contactId string, from time.Time, to time.Time, types []model.ActionType) (*entity.ActionEntities, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	var nodeLabels = []string{}
	for _, v := range types {
		nodeLabels = append(nodeLabels, entity.NodeLabelsByActionType[v.String()])
	}

	dbNodes, err := s.repositories.ActionRepository.GetContactActions(ctx, session, common.GetContext(ctx).Tenant, contactId, from, to, nodeLabels)
	if err != nil {
		return nil, err
	}

	actions := entity.ActionEntities{}
	for _, v := range dbNodes {
		if slices.Contains(v.Labels, entity.NodeLabel_PageView) {
			actions = append(actions, s.mapDbNodeToPageViewAction(v))
		} else if slices.Contains(v.Labels, entity.NodeLabel_InteractionSession) {
			actions = append(actions, s.services.InteractionSessionService.mapDbNodeToInteractionSessionEntity(*v))
		} else if slices.Contains(v.Labels, entity.NodeLabel_Ticket) {
			actions = append(actions, s.services.TicketService.mapDbNodeToTicket(*v))
		} else if slices.Contains(v.Labels, entity.NodeLabel_Conversation) {
			actions = append(actions, s.services.ConversationService.mapDbNodeToConversationEntity(*v))
		} else if slices.Contains(v.Labels, entity.NodeLabel_Note) {
			actions = append(actions, s.services.NoteService.mapDbNodeToNoteEntity(*v))
		}
	}

	return &actions, nil
}

func (s *actionsService) mapDbNodeToPageViewAction(node *dbtype.Node) *entity.PageViewEntity {
	props := utils.GetPropsFromNode(*node)
	pageViewAction := entity.PageViewEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		Application:    utils.GetStringPropOrEmpty(props, "application"),
		TrackerName:    utils.GetStringPropOrEmpty(props, "trackerName"),
		SessionId:      utils.GetStringPropOrEmpty(props, "sessionId"),
		PageUrl:        utils.GetStringPropOrEmpty(props, "pageUrl"),
		PageTitle:      utils.GetStringPropOrEmpty(props, "pageTitle"),
		OrderInSession: utils.GetInt64PropOrZero(props, "orderInSession"),
		EngagedTime:    utils.GetInt64PropOrZero(props, "engagedTime"),
		StartedAt:      utils.GetTimePropOrNow(props, "startedAt"),
		EndedAt:        utils.GetTimePropOrNow(props, "endedAt"),
	}
	return &pageViewAction
}
