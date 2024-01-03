package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/exp/slices"
	"time"
)

type TimelineEventService interface {
	GetTimelineEventsForContact(ctx context.Context, contactId string, from *time.Time, size int, types []model.TimelineEventType) (*entity.TimelineEventEntities, error)
	GetTimelineEventsTotalCountForContact(ctx context.Context, contactId string, types []model.TimelineEventType) (int64, error)
	GetTimelineEventsForOrganization(ctx context.Context, organizationId string, from *time.Time, size int, types []model.TimelineEventType) (*entity.TimelineEventEntities, error)
	GetTimelineEventsTotalCountForOrganization(ctx context.Context, organizationId string, types []model.TimelineEventType) (int64, error)
	GetTimelineEventsWithIds(ctx context.Context, ids []string) (*entity.TimelineEventEntities, error)
}

type timelineEventService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewTimelineEventService(log logger.Logger, repositories *repository.Repositories, services *Services) TimelineEventService {
	return &timelineEventService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *timelineEventService) GetTimelineEventsForContact(ctx context.Context, contactId string, from *time.Time, size int, types []model.TimelineEventType) (*entity.TimelineEventEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventService.GetTimelineEventsForContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.Int("size", size), log.Object("types", types))
	if from != nil {
		span.LogFields(log.String("from", from.String()))
	}

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventService.GetTimelineEventsForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId), log.Int("size", size), log.Object("types", types))
	if from != nil {
		span.LogFields(log.String("from", from.String()))
	}

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventService.GetTimelineEventsTotalCountForContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contactId", contactId), log.Object("types", types))

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventService.GetTimelineEventsTotalCountForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId), log.Object("types", types))

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
	timelineEvents := make(entity.TimelineEventEntities, 0, len(dbNodes))
	for _, v := range dbNodes {
		timelineEvents = append(timelineEvents, s.convertDbNodeToTimelineEvent(v))
	}
	return timelineEvents
}

func (s *timelineEventService) convertDbNodeToTimelineEvent(dbNode *dbtype.Node) entity.TimelineEvent {
	if slices.Contains(dbNode.Labels, neo4jentity.NodeLabel_PageView) {
		return s.services.PageViewService.mapDbNodeToPageView(*dbNode)
	} else if slices.Contains(dbNode.Labels, neo4jentity.NodeLabel_InteractionSession) {
		return s.services.InteractionSessionService.mapDbNodeToInteractionSessionEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, neo4jentity.NodeLabel_Issue) {
		return s.services.IssueService.mapDbNodeToIssue(*dbNode)
	} else if slices.Contains(dbNode.Labels, neo4jentity.NodeLabel_Note) {
		return s.services.NoteService.mapDbNodeToNoteEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, neo4jentity.NodeLabel_InteractionEvent) {
		return s.services.InteractionEventService.mapDbNodeToInteractionEventEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, neo4jentity.NodeLabel_Analysis) {
		return s.services.AnalysisService.mapDbNodeToAnalysisEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, neo4jentity.NodeLabel_Meeting) {
		return s.services.MeetingService.mapDbNodeToMeetingEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, neo4jentity.NodeLabel_Action) {
		return s.services.ActionService.mapDbNodeToActionEntity(*dbNode)
	} else if slices.Contains(dbNode.Labels, neo4jentity.NodeLabel_LogEntry) {
		return s.services.LogEntryService.mapDbNodeToLogEntryEntity(dbNode)
	}
	return nil
}

func (s *timelineEventService) GetTimelineEventsWithIds(ctx context.Context, ids []string) (*entity.TimelineEventEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TimelineEventService.GetTimelineEventsWithIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	dbNodes, err := s.repositories.TimelineEventRepository.GetTimelineEventsWithIds(ctx, common.GetTenantFromContext(ctx), ids)
	if err != nil {
		return nil, err
	}

	timelineEvents := make(entity.TimelineEventEntities, 0, len(dbNodes))
	for _, v := range dbNodes {
		timelineEvent := s.convertDbNodeToTimelineEvent(v)
		timelineEvent.SetDataloaderKey(utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*v), "id"))
		timelineEvents = append(timelineEvents, timelineEvent)
	}

	return &timelineEvents, nil
}
