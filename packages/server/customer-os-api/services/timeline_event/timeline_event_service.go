package api_timeline_event

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	model2 "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"golang.org/x/exp/slices"

	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
)

type timelineEventService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewTimelineEventService(log logger.Logger, repositories *repository.Repositories) cosapi_interfaces.TimelineEventService {
	return &timelineEventService{
		log:          log,
		repositories: repositories,
	}
}

func (s *timelineEventService) GetTimelineEventsForContact(ctx context.Context, contactId string, from *time.Time, size int, types []model.TimelineEventType) (*entity.TimelineEventEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TimelineEventService.GetTimelineEventsForContact")
	defer spans.Finish()
	spans.LogKV("contactId", contactId, "size", size, "types", types)
	if from != nil {
		spans.LogKV("from", from.String())
	}

	nodeLabels := []string{}
	for _, v := range types {
		nodeLabels = append(nodeLabels, entity.NodeLabelsByTimelineEventType[v.String()])
	}
	if len(nodeLabels) == 0 {
		for _, v := range entity.NodeLabelsByTimelineEventType {
			nodeLabels = append(nodeLabels, v)
		}
	}

	var startingDate time.Time
	if from == nil {
		startingDate = utils.Now().Add(time.Duration(5) * time.Second)
	} else {
		startingDate = *from
	}

	dbNodes, err := s.repositories.Neo4jRepositories.TimelineEventReadRepository.GetTimelineEventsForContact(ctx, common.GetContext(ctx).Tenant, contactId, startingDate, size, nodeLabels)
	if err != nil {
		return nil, err
	}

	timelineEvents := s.convertDbNodesToTimelineEvents(dbNodes)

	return &timelineEvents, nil
}

func (s *timelineEventService) GetTimelineEventsForOrganization(ctx context.Context, organizationId string, from *time.Time, size int, types []model.TimelineEventType) (*entity.TimelineEventEntities, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "TimelineEventService.GetTimelineEventsForOrganization")
	defer span.Finish()
	span.LogKV("organizationId", organizationId, "size", size, "types", types)
	if from != nil {
		span.LogKV("from", from.String())
	}

	var nodeLabels []string
	for _, v := range types {
		nodeLabels = append(nodeLabels, entity.NodeLabelsByTimelineEventType[v.String()])
	}
	if len(nodeLabels) == 0 {
		for _, v := range entity.NodeLabelsByTimelineEventType {
			nodeLabels = append(nodeLabels, v)
		}
	}

	var startingDate time.Time
	if from == nil {
		startingDate = utils.Now().Add(time.Duration(5) * time.Second)
	} else {
		startingDate = *from
	}

	dbNodes, err := s.repositories.Neo4jRepositories.TimelineEventReadRepository.GetTimelineEventsForOrganization(ctx, common.GetContext(ctx).Tenant, organizationId, startingDate, size, nodeLabels)
	if err != nil {
		return nil, err
	}

	timelineEvents := s.convertDbNodesToTimelineEvents(dbNodes)
	return &timelineEvents, nil
}

func (s *timelineEventService) GetTimelineEventsTotalCountForContact(ctx context.Context, contactId string, types []model.TimelineEventType) (int64, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "TimelineEventService.GetTimelineEventsTotalCountForContact")
	defer span.Finish()
	span.LogKV("contactId", contactId, "types", types)

	nodeLabels := []string{}
	for _, v := range types {
		nodeLabels = append(nodeLabels, entity.NodeLabelsByTimelineEventType[v.String()])
	}
	if len(nodeLabels) == 0 {
		for _, v := range entity.NodeLabelsByTimelineEventType {
			nodeLabels = append(nodeLabels, v)
		}
	}

	count, err := s.repositories.Neo4jRepositories.TimelineEventReadRepository.GetTimelineEventsTotalCountForContact(ctx, common.GetContext(ctx).Tenant, contactId, nodeLabels)
	if err != nil {
		return int64(0), err
	}

	return count, nil
}

func (s *timelineEventService) GetTimelineEventsTotalCountForOrganization(ctx context.Context, organizationId string, types []model.TimelineEventType) (int64, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "TimelineEventService.GetTimelineEventsTotalCountForOrganization")
	defer span.Finish()
	span.LogKV("organizationId", organizationId, "types", types)

	var nodeLabels []string
	for _, value := range types {
		nodeLabels = append(nodeLabels, entity.NodeLabelsByTimelineEventType[value.String()])
	}
	if len(nodeLabels) == 0 {
		for _, v := range entity.NodeLabelsByTimelineEventType {
			nodeLabels = append(nodeLabels, v)
		}
	}

	count, err := s.repositories.Neo4jRepositories.TimelineEventReadRepository.GetTimelineEventsTotalCountForOrganization(ctx, common.GetContext(ctx).Tenant, organizationId, nodeLabels)
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
	if slices.Contains(dbNode.Labels, model2.NodeLabelInteractionEvent) {
		return neo4jmapper.MapDbNodeToInteractionEventEntity(dbNode)
	} else if slices.Contains(dbNode.Labels, model2.NodeLabelAction) {
		return neo4jmapper.MapDbNodeToActionEntity(dbNode)
	} else if slices.Contains(dbNode.Labels, model2.NodeLabelLogEntry) {
		return neo4jmapper.MapDbNodeToLogEntryEntity(dbNode)
	} else if slices.Contains(dbNode.Labels, model2.NodeLabelMarkdownEvent) {
		return neo4jmapper.MapDbNodeToMarkdownEventEntity(dbNode)
	} else if slices.Contains(dbNode.Labels, model2.NodeLabelIssue) {
		return neo4jmapper.MapDbNodeToIssueEntity(dbNode)
	} else if slices.Contains(dbNode.Labels, model2.NodeLabelMeeting) {
		return neo4jmapper.MapDbNodeToMeetingEntity(dbNode)
	}
	return nil
}

func (s *timelineEventService) GetTimelineEventsWithIds(ctx context.Context, ids []string) (*entity.TimelineEventEntities, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "TimelineEventService.GetTimelineEventsWithIds")
	defer span.Finish()
	span.LogObjectAsJson("ids", ids)

	dbNodes, err := s.repositories.Neo4jRepositories.TimelineEventReadRepository.GetTimelineEventsWithIds(ctx, common.GetTenantFromContext(ctx), ids)
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

func (s *timelineEventService) GetInboundCommsCountCountByOrganizations(ctx context.Context, organizationIds []string) (map[string]int64, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TimelineEventService.GetInboundCommsCountCountByOrganizations")
	defer spans.Finish()
	spans.LogObjectAsJson("organizationIds", organizationIds)

	return s.repositories.Neo4jRepositories.TimelineEventReadRepository.GetInboundCommsTimelineEventsCountByOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
}

func (s *timelineEventService) GetOutboundCommsCountCountByOrganizations(ctx context.Context, organizationIds []string) (map[string]int64, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TimelineEventService.GetOutboundCommsCountCountByOrganizations")
	defer spans.Finish()
	spans.LogObjectAsJson("organizationIds", organizationIds)

	return s.repositories.Neo4jRepositories.TimelineEventReadRepository.GetOutboundCommsTimelineEventsCountByOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
}
