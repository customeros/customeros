package graph_low_prio

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type OrganizationEventHandler struct {
	repositories *repository.Repositories
	log          logger.Logger
	grpcClients  *grpc_client.Clients
}

func NewOrganizationEventHandler(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) *OrganizationEventHandler {
	return &OrganizationEventHandler{
		repositories: repositories,
		log:          log,
		grpcClients:  grpcClients,
	}
}

func (h *OrganizationEventHandler) OnRefreshLastTouchPointV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRefreshLastTouchPointV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRefreshLastTouchpointEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	//fetch the real touchpoint
	//if it doesn't exist, check for the Created Action
	var lastTouchpointId string
	var lastTouchpointAt *time.Time
	var timelineEventNode *dbtype.Node
	var err error

	lastTouchpointAt, lastTouchpointId, err = h.repositories.Neo4jRepositories.TimelineEventReadRepository.CalculateAndGetLastTouchPoint(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to calculate last touchpoint: %v", err.Error())
		span.LogFields(log.Bool("last touchpoint failed", true))
		return nil
	}

	if lastTouchpointAt == nil {
		timelineEventNode, err = h.repositories.Neo4jRepositories.ActionReadRepository.GetSingleAction(ctx, eventData.Tenant, organizationId, neo4jenum.ORGANIZATION, neo4jenum.ActionCreated)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to get created action: %v", err.Error())
			return nil
		}
		if timelineEventNode != nil {
			propsFromNode := utils.GetPropsFromNode(*timelineEventNode)
			lastTouchpointId = utils.GetStringPropOrEmpty(propsFromNode, "id")
			lastTouchpointAt = utils.GetTimePropOrNil(propsFromNode, "createdAt")
		}
	} else {
		timelineEventNode, err = h.repositories.Neo4jRepositories.TimelineEventReadRepository.GetTimelineEvent(ctx, eventData.Tenant, lastTouchpointId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to get last touchpoint: %v", err.Error())
			return nil
		}
	}

	if timelineEventNode == nil {
		h.log.Infof("Last touchpoint not available for organization: %s", organizationId)
		span.LogFields(log.Bool("last touchpoint not found", true))
		return nil
	}

	timelineEvent := neo4jmapper.MapDbNodeToTimelineEvent(timelineEventNode)
	if timelineEvent == nil {
		h.log.Infof("Last touchpoint not available for organization: %s", organizationId)
		span.LogFields(log.Bool("last touchpoint not found", true))
		return nil
	}

	var timelineEventType string
	switch timelineEvent.TimelineEventLabel() {
	case neo4jutil.NodeLabelPageView:
		timelineEventType = "PAGE_VIEW"
	case neo4jutil.NodeLabelInteractionSession:
		timelineEventType = "INTERACTION_SESSION"
	case neo4jutil.NodeLabelNote:
		timelineEventType = "NOTE"
	case neo4jutil.NodeLabelInteractionEvent:
		timelineEventInteractionEvent := timelineEvent.(*neo4jentity.InteractionEventEntity)
		if timelineEventInteractionEvent.Channel == "EMAIL" {
			timelineEventType = "INTERACTION_EVENT_EMAIL_SENT"
		} else if timelineEventInteractionEvent.Channel == "VOICE" {
			timelineEventType = "INTERACTION_EVENT_PHONE_CALL"
		} else if timelineEventInteractionEvent.Channel == "CHAT" {
			timelineEventType = "INTERACTION_EVENT_CHAT"
		} else if timelineEventInteractionEvent.EventType == "meeting" {
			timelineEventType = "MEETING"
		}
	case neo4jutil.NodeLabelAnalysis:
		timelineEventType = "ANALYSIS"
	case neo4jutil.NodeLabelMeeting:
		timelineEventType = "MEETING"
	case neo4jutil.NodeLabelAction:
		timelineEventAction := timelineEvent.(*neo4jentity.ActionEntity)
		if timelineEventAction.Type == neo4jenum.ActionCreated {
			timelineEventType = "ACTION_CREATED"
		} else {
			timelineEventType = "ACTION"
		}
	case neo4jutil.NodeLabelLogEntry:
		timelineEventType = "LOG_ENTRY"
	case neo4jutil.NodeLabelIssue:
		timelineEventIssue := timelineEvent.(*neo4jentity.IssueEntity)
		if timelineEventIssue.CreatedAt.Equal(timelineEventIssue.UpdatedAt) {
			timelineEventType = "ISSUE_CREATED"
		} else {
			timelineEventType = "ISSUE_UPDATED"
		}
	default:
		h.log.Infof("Last touchpoint not available for organization: %s", organizationId)
	}

	if err = h.repositories.Neo4jRepositories.OrganizationWriteRepository.UpdateLastTouchpoint(ctx, eventData.Tenant, organizationId, *lastTouchpointAt, lastTouchpointId, timelineEventType); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update last touchpoint for tenant %s, organization %s: %s", eventData.Tenant, organizationId, err.Error())
	}

	return nil
}
