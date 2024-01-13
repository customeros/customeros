package graph

import (
	"context"
	"encoding/json"
	"fmt"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
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

type eventMetadata struct {
	UserId string `json:"user-id"`
}

func (h *OrganizationEventHandler) OnOrganizationCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationCreate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	session := utils.NewNeo4jWriteSession(ctx, *h.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var err error
		data := neo4jrepository.OrganizationCreateFields{
			CreatedAt: eventData.CreatedAt,
			UpdatedAt: eventData.UpdatedAt,
			SourceFields: neo4jmodel.Source{
				Source:        helper.GetSource(eventData.Source),
				SourceOfTruth: helper.GetSource(eventData.SourceOfTruth),
				AppSource:     helper.GetSource(eventData.AppSource),
			},
			Name:               eventData.Name,
			Hide:               eventData.Hide,
			Description:        eventData.Description,
			Website:            eventData.Website,
			Industry:           eventData.Industry,
			SubIndustry:        eventData.SubIndustry,
			IndustryGroup:      eventData.IndustryGroup,
			TargetAudience:     eventData.TargetAudience,
			ValueProposition:   eventData.ValueProposition,
			IsPublic:           eventData.IsPublic,
			IsCustomer:         eventData.IsCustomer,
			Employees:          eventData.Employees,
			Market:             eventData.Market,
			LastFundingRound:   eventData.LastFundingRound,
			LastFundingAmount:  eventData.LastFundingAmount,
			ReferenceId:        eventData.ReferenceId,
			Note:               eventData.Note,
			LogoUrl:            eventData.LogoUrl,
			Headquarters:       eventData.Headquarters,
			YearFounded:        eventData.YearFounded,
			EmployeeGrowthRate: eventData.EmployeeGrowthRate,
		}
		err = h.repositories.Neo4jRepositories.OrganizationWriteRepository.CreateOrganizationInTx(ctx, tx, eventData.Tenant, organizationId, data)
		if err != nil {
			h.log.Errorf("Error while saving organization %s: %s", organizationId, err.Error())
			return nil, err
		}
		if eventData.ExternalSystem.Available() {
			externalSystemData := neo4jmodel.ExternalSystem{
				ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
				ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
				ExternalId:       eventData.ExternalSystem.ExternalId,
				ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
				ExternalSource:   eventData.ExternalSystem.ExternalSource,
				SyncDate:         eventData.ExternalSystem.SyncDate,
			}
			err = h.repositories.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, organizationId, neo4jutil.NodeLabelOrganization, externalSystemData)
			if err != nil {
				h.log.Errorf("Error while link organization %s with external system %s: %s", organizationId, eventData.ExternalSystem.ExternalSystemId, err.Error())
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// set customer os id
	customerOsErr := h.setCustomerOsId(ctx, eventData.Tenant, organizationId)
	if customerOsErr != nil {
		tracing.TraceErr(span, customerOsErr)
		h.log.Errorf("Failed to set customer os id for tenant %s organization %s", eventData.Tenant, organizationId)
	}

	// Set organization owner
	evtMetadata := eventMetadata{}
	if err = json.Unmarshal(evt.Metadata, &evtMetadata); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Unmarshal")
	} else {
		if evtMetadata.UserId != "" {
			err = h.repositories.Neo4jRepositories.OrganizationWriteRepository.ReplaceOwner(ctx, eventData.Tenant, organizationId, evtMetadata.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to replace owner of organization %s with user %s", organizationId, evtMetadata.UserId)
			}
		}
	}

	// Set create action
	_, err = h.repositories.Neo4jRepositories.ActionWriteRepository.MergeByActionType(ctx, eventData.Tenant, organizationId, neo4jenum.ORGANIZATION, neo4jenum.ActionCreated, "", "", eventData.CreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed creating likelihood update action for organization %s: %s", organizationId, err.Error())
	}

	// Request last touch point update
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = h.grpcClients.OrganizationClient.RefreshLastTouchpoint(ctx, &organizationpb.OrganizationIdGrpcRequest{
		Tenant:         eventData.Tenant,
		OrganizationId: organizationId,
		AppSource:      constants.AppSourceEventProcessingPlatform,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while refreshing last touchpoint for organization %s: %s", organizationId, err.Error())
	}

	return nil
}

func (h *OrganizationEventHandler) setCustomerOsId(ctx context.Context, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.setCustomerOsId")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("OrganizationId", organizationId))

	orgDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)

	if organizationEntity.CustomerOsId != "" {
		return nil
	}
	var customerOsId string
	maxAttempts := 20
	for attempt := 1; attempt < maxAttempts+1; attempt++ {
		customerOsId = generateNewRandomCustomerOsId()
		customerOsIdsEntity := postgresentity.CustomerOsIds{
			Tenant:       tenant,
			CustomerOSID: customerOsId,
			Entity:       postgresentity.Organization,
			EntityId:     organizationId,
			Attempts:     attempt,
		}
		innerErr := h.repositories.CustomerOsIdsRepository.Reserve(customerOsIdsEntity)
		if innerErr == nil {
			break
		}
	}
	return h.repositories.Neo4jRepositories.OrganizationWriteRepository.SetCustomerOsIdIfMissing(ctx, tenant, organizationId, customerOsId)
}

func (h *OrganizationEventHandler) OnOrganizationUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationUpdate")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	data := neo4jrepository.OrganizationUpdateFields{
		UpdatedAt:                eventData.UpdatedAt,
		Name:                     eventData.Name,
		Hide:                     eventData.Hide,
		Description:              eventData.Description,
		Website:                  eventData.Website,
		Industry:                 eventData.Industry,
		SubIndustry:              eventData.SubIndustry,
		IndustryGroup:            eventData.IndustryGroup,
		TargetAudience:           eventData.TargetAudience,
		ValueProposition:         eventData.ValueProposition,
		IsPublic:                 eventData.IsPublic,
		IsCustomer:               eventData.IsCustomer,
		Employees:                eventData.Employees,
		Market:                   eventData.Market,
		LastFundingRound:         eventData.LastFundingRound,
		LastFundingAmount:        eventData.LastFundingAmount,
		ReferenceId:              eventData.ReferenceId,
		Note:                     eventData.Note,
		LogoUrl:                  eventData.LogoUrl,
		Headquarters:             eventData.Headquarters,
		YearFounded:              eventData.YearFounded,
		EmployeeGrowthRate:       eventData.EmployeeGrowthRate,
		WebScrapedUrl:            eventData.WebScrapedUrl,
		Source:                   helper.GetSource(eventData.Source),
		UpdateName:               eventData.UpdateName(),
		UpdateDescription:        eventData.UpdateDescription(),
		UpdateHide:               eventData.UpdateHide(),
		UpdateIsCustomer:         eventData.UpdateIsCustomer(),
		UpdateWebsite:            eventData.UpdateWebsite(),
		UpdateIndustry:           eventData.UpdateIndustry(),
		UpdateSubIndustry:        eventData.UpdateSubIndustry(),
		UpdateIndustryGroup:      eventData.UpdateIndustryGroup(),
		UpdateTargetAudience:     eventData.UpdateTargetAudience(),
		UpdateValueProposition:   eventData.UpdateValueProposition(),
		UpdateLastFundingRound:   eventData.UpdateLastFundingRound(),
		UpdateLastFundingAmount:  eventData.UpdateLastFundingAmount(),
		UpdateReferenceId:        eventData.UpdateReferenceId(),
		UpdateNote:               eventData.UpdateNote(),
		UpdateIsPublic:           eventData.UpdateIsPublic(),
		UpdateEmployees:          eventData.UpdateEmployees(),
		UpdateMarket:             eventData.UpdateMarket(),
		UpdateYearFounded:        eventData.UpdateYearFounded(),
		UpdateHeadquarters:       eventData.UpdateHeadquarters(),
		UpdateLogoUrl:            eventData.UpdateLogoUrl(),
		UpdateEmployeeGrowthRate: eventData.UpdateEmployeeGrowthRate(),
	}
	err := h.repositories.Neo4jRepositories.OrganizationWriteRepository.UpdateOrganization(ctx, eventData.Tenant, organizationId, data)
	// set customer os id
	customerOsErr := h.setCustomerOsId(ctx, eventData.Tenant, organizationId)
	if customerOsErr != nil {
		tracing.TraceErr(span, customerOsErr)
		h.log.Errorf("Failed to set customer os id for tenant %s organization %s", eventData.Tenant, organizationId)
	}
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving organization %s: %s", organizationId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		session := utils.NewNeo4jWriteSession(ctx, *h.repositories.Drivers.Neo4jDriver)
		defer session.Close(ctx)

		_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			if eventData.ExternalSystem.Available() {
				externalSystemData := neo4jmodel.ExternalSystem{
					ExternalSystemId: eventData.ExternalSystem.ExternalSystemId,
					ExternalUrl:      eventData.ExternalSystem.ExternalUrl,
					ExternalId:       eventData.ExternalSystem.ExternalId,
					ExternalIdSecond: eventData.ExternalSystem.ExternalIdSecond,
					ExternalSource:   eventData.ExternalSystem.ExternalSource,
					SyncDate:         eventData.ExternalSystem.SyncDate,
				}
				innerErr := h.repositories.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, organizationId, neo4jutil.NodeLabelOrganization, externalSystemData)
				if innerErr != nil {
					h.log.Errorf("Error while link organization %s with external system %s: %s", organizationId, eventData.ExternalSystem.ExternalSystemId, err.Error())
					return nil, innerErr
				}
			}
			return nil, nil
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (h *OrganizationEventHandler) OnPhoneNumberLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnPhoneNumberLinkedToOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.PhoneNumberId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *OrganizationEventHandler) OnEmailLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnEmailLinkedToOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.EmailWriteRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.EmailId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *OrganizationEventHandler) OnLocationLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnLocationLinkedToOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationLinkLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.LocationWriteRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.LocationId, eventData.UpdatedAt)

	return err
}

func (h *OrganizationEventHandler) OnDomainLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnDomainLinkedToOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationLinkDomainEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	if strings.TrimSpace(eventData.Domain) == "" {
		return nil
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	if !utils.IsValidTLD(eventData.Domain) {
		err := errors.New(fmt.Sprintf("Invalid domain: %s", eventData.Domain))
		err = errors.Wrap(err, "IsValidTLD")
		tracing.TraceErr(span, err)
		h.log.Error("Not linked domain to organization %s : %s", organizationId, err.Error())
		return nil
	}

	err := h.repositories.Neo4jRepositories.OrganizationWriteRepository.LinkWithDomain(ctx, eventData.Tenant, organizationId, strings.TrimSpace(eventData.Domain))

	return err
}

func (h *OrganizationEventHandler) OnSocialAddedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnSocialAddedToOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationAddSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	data := neo4jrepository.SocialFields{
		SocialId:     eventData.SocialId,
		Url:          eventData.Url,
		PlatformName: eventData.PlatformName,
		CreatedAt:    eventData.CreatedAt,
		UpdatedAt:    eventData.UpdatedAt,
		SourceFields: neo4jmodel.Source{
			Source:        helper.GetSource(eventData.Source),
			SourceOfTruth: helper.GetSource(eventData.SourceOfTruth),
			AppSource:     helper.GetSource(eventData.AppSource),
		},
	}
	err := h.repositories.Neo4jRepositories.SocialWriteRepository.MergeSocialFor(ctx, eventData.Tenant, organizationId, neo4jutil.NodeLabelOrganization, data)

	return err
}

func (h *OrganizationEventHandler) OnOrganizationHide(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationHide")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.HideOrganizationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.OrganizationWriteRepository.SetVisibility(ctx, eventData.Tenant, organizationId, true)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (h *OrganizationEventHandler) OnOrganizationShow(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationShow")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.HideOrganizationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.Neo4jRepositories.OrganizationWriteRepository.SetVisibility(ctx, eventData.Tenant, organizationId, false)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	// set customer os id
	customerOsErr := h.setCustomerOsId(ctx, eventData.Tenant, organizationId)
	if customerOsErr != nil {
		tracing.TraceErr(span, customerOsErr)
		h.log.Errorf("Failed to set customer os id for tenant %s organization %s", eventData.Tenant, organizationId)
	}

	return err
}

func (h *OrganizationEventHandler) OnRefreshLastTouchpoint(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRefreshLastTouchpoint")
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

	timelineEvent := graph_db.MapDbNodeToTimelineEvent(timelineEventNode)
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
		timelineEventInteractionEvent := timelineEvent.(*entity.InteractionEventEntity)
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
		timelineEventAction := timelineEvent.(*entity.ActionEntity)
		if timelineEventAction.Type == neo4jenum.ActionCreated {
			timelineEventType = "ACTION_CREATED"
		} else {
			timelineEventType = "ACTION"
		}
	case neo4jutil.NodeLabelLogEntry:
		timelineEventType = "LOG_ENTRY"
	case neo4jutil.NodeLabelIssue:
		timelineEventIssue := timelineEvent.(*entity.IssueEntity)
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

func (h *OrganizationEventHandler) OnRefreshArr(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRefreshArr")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRefreshArrEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	if err := h.repositories.Neo4jRepositories.OrganizationWriteRepository.UpdateArr(ctx, eventData.Tenant, organizationId); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update arr for tenant %s, organization %s: %s", eventData.Tenant, organizationId, err.Error())
	}

	return nil
}

func (h *OrganizationEventHandler) OnRefreshRenewalSummary(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRefreshRenewalSummary")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRefreshArrEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	openRenewalOpportunityDbNodes, err := h.repositories.Neo4jRepositories.OpportunityReadRepository.GetOpenRenewalOpportunitiesForOrganization(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to get open renewal opportunities for organization %s: %s", organizationId, err.Error())
		return nil
	}
	var nextRenewalDate *time.Time
	var lowestRenewalLikelihood *string
	var renewalLikelihoodOrder int64
	if len(openRenewalOpportunityDbNodes) > 0 {
		opportunities := make([]entity.OpportunityEntity, len(openRenewalOpportunityDbNodes))
		for _, opportunityDbNode := range openRenewalOpportunityDbNodes {
			opportunities = append(opportunities, *graph_db.MapDbNodeToOpportunityEntity(opportunityDbNode))
		}
		for _, opportunity := range opportunities {
			if opportunity.RenewalDetails.RenewedAt != nil && opportunity.RenewalDetails.RenewedAt.After(utils.Now()) {
				if nextRenewalDate == nil || opportunity.RenewalDetails.RenewedAt.Before(*nextRenewalDate) {
					nextRenewalDate = opportunity.RenewalDetails.RenewedAt
				}
			}
			if opportunity.RenewalDetails.RenewalLikelihood != "" {
				order := getOrderForRenewalLikelihood(opportunity.RenewalDetails.RenewalLikelihood)
				if renewalLikelihoodOrder == 0 || renewalLikelihoodOrder > order {
					renewalLikelihoodOrder = order
					lowestRenewalLikelihood = utils.ToPtr(opportunity.RenewalDetails.RenewalLikelihood)
				}
			}
		}
	}

	renewalLikelihoodOrderPtr := utils.ToPtr[int64](renewalLikelihoodOrder)
	if renewalLikelihoodOrder == 0 {
		renewalLikelihoodOrderPtr = nil
	}

	if err := h.repositories.Neo4jRepositories.OrganizationWriteRepository.UpdateRenewalSummary(ctx, eventData.Tenant, organizationId, lowestRenewalLikelihood, renewalLikelihoodOrderPtr, nextRenewalDate); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update arr for tenant %s, organization %s: %s", eventData.Tenant, organizationId, err.Error())
	}

	return nil
}

func getOrderForRenewalLikelihood(likelihood string) int64 {
	switch likelihood {
	case string(neo4jenum.RenewalLikelihoodHigh):
		return constants.RenewalLikelihood_Order_High
	case string(neo4jenum.RenewalLikelihoodMedium):
		return constants.RenewalLikelihood_Order_Medium
	case string(neo4jenum.RenewalLikelihoodLow):
		return constants.RenewalLikelihood_Order_Low
	case string(neo4jenum.RenewalLikelihoodZero):
		return constants.RenewalLikelihood_Order_Zero
	default:
		return 0
	}
}

func (h *OrganizationEventHandler) OnUpsertCustomField(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnUpsertCustomField")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationUpsertCustomField
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	customFieldExists, err := h.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, eventData.Tenant, eventData.CustomFieldId, neo4jutil.NodeLabelCustomField)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to check if custom field exists: %s", err.Error())
		return err
	}
	if !customFieldExists {
		data := neo4jrepository.CustomFieldCreateFields{
			CreatedAt:           eventData.CreatedAt,
			UpdatedAt:           eventData.UpdatedAt,
			ExistsInEventStore:  eventData.ExistsInEventStore,
			TemplateId:          eventData.TemplateId,
			CustomFieldId:       eventData.CustomFieldId,
			CustomFieldName:     eventData.CustomFieldName,
			CustomFieldDataType: eventData.CustomFieldDataType,
			CustomFieldValue:    eventData.CustomFieldValue,
			SourceFields: neo4jmodel.Source{
				Source:        helper.GetSource(eventData.Source),
				SourceOfTruth: helper.GetSource(eventData.SourceOfTruth),
				AppSource:     helper.GetSource(eventData.AppSource),
			},
		}
		err = h.repositories.Neo4jRepositories.CustomFieldWriteRepository.AddCustomFieldToOrganization(ctx, eventData.Tenant, organizationId, data)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to add custom field to organization: %s", err.Error())
			return err
		}
	} else {
		//TODO implement update custom field
	}

	return nil
}

func (h *OrganizationEventHandler) OnLinkWithParentOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnLinkWithParentOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationAddParentEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	return h.repositories.Neo4jRepositories.OrganizationWriteRepository.LinkWithParentOrganization(ctx, eventData.Tenant, organizationId, eventData.ParentOrganizationId, eventData.Type)
}

func (h *OrganizationEventHandler) OnUnlinkFromParentOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnUnlinkFromParentOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRemoveParentEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	return h.repositories.Neo4jRepositories.OrganizationWriteRepository.UnlinkParentOrganization(ctx, eventData.Tenant, organizationId, eventData.ParentOrganizationId)
}

func (h *OrganizationEventHandler) OnUpdateOnboardingStatus(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnUpdateOnboardingStatus")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UpdateOnboardingStatusEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	organizationDbNode, err := h.repositories.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to get organization %s: %s", organizationId, err.Error())
		return err
	}
	if organizationDbNode == nil {
		err = errors.New(fmt.Sprintf("Organization %s not found", organizationId))
		tracing.TraceErr(span, err)
		return nil
	}
	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	err = h.repositories.Neo4jRepositories.OrganizationWriteRepository.UpdateOnboardingStatus(ctx, eventData.Tenant, organizationId, eventData.Status, eventData.Comments, getOrderForOnboardingStatus(eventData.Status), eventData.UpdatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update onboarding status for organization %s: %s", organizationId, err.Error())
		return err
	}

	if eventData.CausedByContractId != "" {
		err = h.repositories.Neo4jRepositories.ContractWriteRepository.ContractCausedOnboardingStatusChange(ctx, eventData.Tenant, eventData.CausedByContractId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to update contract %s caused onboarding status change: %s", eventData.CausedByContractId, err.Error())
		}
	}

	if organizationEntity.OnboardingDetails.Status != eventData.Status {
		err = h.saveOnboardingStatusChangeAction(ctx, organizationId, eventData, span)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to save onboarding status change action for organization %s: %s", organizationId, err.Error())
		}
	}

	return nil
}

func getOrderForOnboardingStatus(status string) *int64 {
	switch status {
	case string(model.OnboardingStatusNotStarted):
		return utils.Int64Ptr(constants.OnboardingStatus_Order_NotStarted)
	case string(model.OnboardingStatusOnTrack):
		return utils.Int64Ptr(constants.OnboardingStatus_Order_OnTrack)
	case string(model.OnboardingStatusLate):
		return utils.Int64Ptr(constants.OnboardingStatus_Order_Late)
	case string(model.OnboardingStatusStuck):
		return utils.Int64Ptr(constants.OnboardingStatus_Order_Stuck)
	case string(model.OnboardingStatusDone):
		return utils.Int64Ptr(constants.OnboardingStatus_Order_Done)
	case string(model.OnboardingStatusSuccessful):
		return utils.Int64Ptr(constants.OnboardingStatus_Order_Successful)
	default:
		return nil
	}
}

type ActionOnboardingStatusMetadata struct {
	Status     string `json:"status"`
	Comments   string `json:"comments"`
	UserId     string `json:"userId"`
	ContractId string `json:"contractId"`
}

func (h *OrganizationEventHandler) saveOnboardingStatusChangeAction(ctx context.Context, organizationId string, eventData events.UpdateOnboardingStatusEvent, span opentracing.Span) error {
	metadata, _ := utils.ToJson(ActionOnboardingStatusMetadata{
		Status:     eventData.Status,
		Comments:   eventData.Comments,
		UserId:     eventData.UpdatedByUserId,
		ContractId: eventData.CausedByContractId,
	})
	message := ""
	userName := ""
	if eventData.UpdatedByUserId != "" {
		userDbNode, err := h.repositories.Neo4jRepositories.UserReadRepository.GetUserById(ctx, eventData.Tenant, eventData.UpdatedByUserId)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to get user %s: %s", eventData.UpdatedByUserId, err.Error())
		}
		if userDbNode != nil {
			user := neo4jmapper.MapDbNodeToUserEntity(userDbNode)
			userName = user.GetFullName()
		}
	}
	if eventData.UpdatedByUserId != "" {
		message = fmt.Sprintf("%s changed the onboarding status to %s", userName, onboardingStatusReadableStringForActionMessage(eventData.Status))
	} else {
		message = fmt.Sprintf("The onboarding status was automatically set to %s", onboardingStatusReadableStringForActionMessage(eventData.Status))
	}

	extraActionProperties := map[string]interface{}{
		"status":   eventData.Status,
		"comments": eventData.Comments,
	}
	_, err := h.repositories.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, organizationId, neo4jenum.ORGANIZATION, neo4jenum.ActionOnboardingStatusChanged, message, metadata, eventData.UpdatedAt, extraActionProperties)
	return err
}

func onboardingStatusReadableStringForActionMessage(status string) string {
	switch status {
	case string(neo4jenum.OnboardingStatusNotApplicable):
		return "Not applicable"
	case string(neo4jenum.OnboardingStatusNotStarted):
		return "Not started"
	case string(neo4jenum.OnboardingStatusOnTrack):
		return "On track"
	case string(neo4jenum.OnboardingStatusLate):
		return "Late"
	case string(neo4jenum.OnboardingStatusStuck):
		return "Stuck"
	case string(neo4jenum.OnboardingStatusDone):
		return "Done"
	case string(neo4jenum.OnboardingStatusSuccessful):
		return "Successful"
	default:
		return status
	}
}

func (h *OrganizationEventHandler) OnUpdateOwner(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnUpdateOrganizationOwner")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationOwnerUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	return h.repositories.Neo4jRepositories.OrganizationWriteRepository.ReplaceOwner(ctx, eventData.Tenant, eventData.OrganizationId, eventData.OwnerUserId)
}

func (h *OrganizationEventHandler) OnCreateBillingProfile(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnCreateBillingProfile")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.BillingProfileCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	data := neo4jrepository.BillingProfileCreateFields{
		OrganizationId: organizationId,
		LegalName:      eventData.LegalName,
		TaxId:          eventData.TaxId,
		CreatedAt:      eventData.CreatedAt,
		UpdatedAt:      eventData.UpdatedAt,
		SourceFields: neo4jmodel.Source{
			Source:    helper.GetSource(eventData.SourceFields.Source),
			AppSource: helper.GetSource(eventData.SourceFields.AppSource),
		},
	}
	return h.repositories.Neo4jRepositories.BillingProfileWriteRepository.Create(ctx, eventData.Tenant, eventData.BillingProfileId, data)
}

func (h *OrganizationEventHandler) OnUpdateBillingProfile(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnUpdateBillingProfile")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.BillingProfileUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	data := neo4jrepository.BillingProfileUpdateFields{
		OrganizationId:  organizationId,
		UpdatedAt:       eventData.UpdatedAt,
		LegalName:       eventData.LegalName,
		TaxId:           eventData.TaxId,
		UpdateLegalName: eventData.UpdateLegalName(),
		UpdateTaxId:     eventData.UpdateTaxId(),
	}
	return h.repositories.Neo4jRepositories.BillingProfileWriteRepository.Update(ctx, eventData.Tenant, eventData.BillingProfileId, data)
}

func (h *OrganizationEventHandler) OnEmailLinkedToBillingProfile(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnEmailLinkedToBillingProfile")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.LinkEmailToBillingProfileEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	return h.repositories.Neo4jRepositories.BillingProfileWriteRepository.LinkEmailToBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.EmailId, eventData.Primary, eventData.UpdatedAt)
}

func (h *OrganizationEventHandler) OnEmailUnlinkedFromBillingProfile(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnEmailUnlinkedFromBillingProfile")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UnlinkEmailFromBillingProfileEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	return h.repositories.Neo4jRepositories.BillingProfileWriteRepository.UnlinkEmailFromBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.EmailId, eventData.UpdatedAt)
}

func (h *OrganizationEventHandler) OnLocationLinkedToBillingProfile(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnLocationLinkedToBillingProfile")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.LinkLocationToBillingProfileEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	return h.repositories.Neo4jRepositories.BillingProfileWriteRepository.LinkLocationToBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.LocationId, eventData.UpdatedAt)
}

func (h *OrganizationEventHandler) OnLocationUnlinkedFromBillingProfile(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnLocationUnlinkedFromBillingProfile")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.UnlinkLocationFromBillingProfileEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	return h.repositories.Neo4jRepositories.BillingProfileWriteRepository.UnlinkLocationFromBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.LocationId, eventData.UpdatedAt)
}
