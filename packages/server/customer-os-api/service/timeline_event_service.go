package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/exp/slices"
	"time"
)

type TimelineEventService interface {
	GetTimelineEventsForContact(ctx context.Context, contactId string, from *time.Time, size int, types []model.TimelineEventType) (*entity.TimelineEventEntities, error)
	GetTimelineEventsTotalCountForContact(ctx context.Context, contactId string, types []model.TimelineEventType) (int64, error)
	GetTimelineEventsForOrganization(ctx context.Context, organizationId string, from *time.Time, size int, types []model.TimelineEventType) (*entity.TimelineEventEntities, error)
	GetTimelineEventsTotalCountForOrganization(ctx context.Context, organizationId string, types []model.TimelineEventType) (int64, error)
}

type timelineEventService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewTimelineEventService(repositories *repository.Repositories, services *Services) TimelineEventService {
	return &timelineEventService{
		repositories: repositories,
		services:     services,
	}
}

func (s *timelineEventService) GetTimelineEventsForContact(ctx context.Context, contactId string, from *time.Time, size int, types []model.TimelineEventType) (*entity.TimelineEventEntities, error) {
	var nodeLabels = []string{}
	for _, v := range types {
		nodeLabels = append(nodeLabels, entity.NodeLabelsByTimelineEventType[v.String()])
	}

	var startingDate time.Time
	if from == nil {
		startingDate = utils.Now().Add(time.Duration(5) * time.Second)
	} else {
		startingDate = *from
	}

	dbNodes, err := s.repositories.TimelineEventRepository.GetTimelineEventsForContact(ctx, common.GetContext(ctx).Tenant, contactId, startingDate, size, nodeLabels)
	if err != nil {
		return nil, err
	}

	timelineEvents := s.convertDbNodesToTimelineEvents(dbNodes)

	return &timelineEvents, nil
}

func (s *timelineEventService) GetTimelineEventsForOrganization(ctx context.Context, organizationId string, from *time.Time, size int, types []model.TimelineEventType) (*entity.TimelineEventEntities, error) {
	var nodeLabels = []string{}
	for _, v := range types {
		nodeLabels = append(nodeLabels, entity.NodeLabelsByTimelineEventType[v.String()])
	}

	var startingDate time.Time
	if from == nil {
		startingDate = utils.Now().Add(time.Duration(5) * time.Second)
	} else {
		startingDate = *from
	}

	dbNodes, err := s.repositories.TimelineEventRepository.GetTimelineEventsForOrganization(ctx, common.GetContext(ctx).Tenant, organizationId, startingDate, size, nodeLabels)
	if err != nil {
		return nil, err
	}

	timelineEvents := s.convertDbNodesToTimelineEvents(dbNodes)
	return &timelineEvents, nil
}

func (s *timelineEventService) GetTimelineEventsTotalCountForContact(ctx context.Context, contactId string, types []model.TimelineEventType) (int64, error) {
	var nodeLabels = []string{}
	for _, v := range types {
		nodeLabels = append(nodeLabels, entity.NodeLabelsByTimelineEventType[v.String()])
	}

	count, err := s.repositories.TimelineEventRepository.GetTimelineEventsTotalCountForContact(ctx, common.GetContext(ctx).Tenant, contactId, nodeLabels)
	if err != nil {
		return int64(0), err
	}

	return count, nil
}

func (s *timelineEventService) GetTimelineEventsTotalCountForOrganization(ctx context.Context, organizationId string, types []model.TimelineEventType) (int64, error) {
	var nodeLabels = []string{}
	for _, v := range types {
		nodeLabels = append(nodeLabels, entity.NodeLabelsByTimelineEventType[v.String()])
	}

	count, err := s.repositories.TimelineEventRepository.GetTimelineEventsTotalCountForOrganization(ctx, common.GetContext(ctx).Tenant, organizationId, nodeLabels)
	if err != nil {
		return int64(0), err
	}

	return count, nil
}

func (s *timelineEventService) convertDbNodesToTimelineEvents(dbNodes []*dbtype.Node) entity.TimelineEventEntities {
	timelineEvents := entity.TimelineEventEntities{}
	for _, v := range dbNodes {
		if slices.Contains(v.Labels, entity.NodeLabel_PageView) {
			timelineEvents = append(timelineEvents, s.services.PageViewService.mapDbNodeToPageView(*v))
		} else if slices.Contains(v.Labels, entity.NodeLabel_InteractionSession) {
			timelineEvents = append(timelineEvents, s.services.InteractionSessionService.mapDbNodeToInteractionSessionEntity(*v))
		} else if slices.Contains(v.Labels, entity.NodeLabel_Ticket) {
			timelineEvents = append(timelineEvents, s.services.TicketService.mapDbNodeToTicket(*v))
		} else if slices.Contains(v.Labels, entity.NodeLabel_Conversation) {
			timelineEvents = append(timelineEvents, s.services.ConversationService.mapDbNodeToConversationEntity(*v))
		} else if slices.Contains(v.Labels, entity.NodeLabel_Note) {
			timelineEvents = append(timelineEvents, s.services.NoteService.mapDbNodeToNoteEntity(*v))
		} else if slices.Contains(v.Labels, entity.NodeLabel_InteractionEvent) {
			timelineEvents = append(timelineEvents, s.services.InteractionEventService.mapDbNodeToInteractionEventEntity(*v))
		} else if slices.Contains(v.Labels, entity.NodeLabel_Analysis) {
			timelineEvents = append(timelineEvents, s.services.AnalysisService.mapDbNodeToAnalysisEntity(*v))
		}
	}
	return timelineEvents
}
