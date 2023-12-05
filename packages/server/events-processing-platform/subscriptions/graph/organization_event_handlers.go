package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	opportunitymodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	cmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type OrganizationEventHandler struct {
	repositories         *repository.Repositories
	organizationCommands *cmdhnd.CommandHandlers
	log                  logger.Logger
}

type eventMetadata struct {
	UserId string `json:"user-id"`
}

func (h *OrganizationEventHandler) OnOrganizationCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationCreate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

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
		err = h.repositories.OrganizationRepository.CreateOrganizationInTx(ctx, tx, organizationId, eventData)
		if err != nil {
			h.log.Errorf("Error while saving organization %s: %s", organizationId, err.Error())
			return nil, err
		}
		if eventData.ExternalSystem.Available() {
			err = h.repositories.ExternalSystemRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, organizationId, constants.NodeLabel_Organization, eventData.ExternalSystem)
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
			err = h.repositories.OrganizationRepository.ReplaceOwner(ctx, eventData.Tenant, organizationId, evtMetadata.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to replace owner of organization %s with user %s", organizationId, evtMetadata.UserId)
			}
		}
	}

	// Set create action
	_, err = h.repositories.ActionRepository.MergeByActionType(ctx, eventData.Tenant, organizationId, entity.ORGANIZATION, entity.ActionCreated, "", "", eventData.CreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed creating likelihood update action for organization %s: %s", organizationId, err.Error())
	}

	// Request last touch point update
	err = h.organizationCommands.RefreshLastTouchpointCommand.Handle(ctx, cmd.NewRefreshLastTouchpointCommand(eventData.Tenant, organizationId, "", constants.AppSourceEventProcessingPlatform))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("RefreshLastTouchpointCommand failed: %v", err.Error())
	}

	return err
}

func (h *OrganizationEventHandler) setCustomerOsId(ctx context.Context, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.setCustomerOsId")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("OrganizationId", organizationId))

	orgDbNode, err := h.repositories.OrganizationRepository.GetOrganization(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	organizationEntity := graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)

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
	return h.repositories.OrganizationRepository.SetCustomerOsIdIfMissing(ctx, tenant, organizationId, customerOsId)
}

func (h *OrganizationEventHandler) OnOrganizationUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationUpdate")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	var err error
	if eventData.IgnoreEmptyFields {
		err = h.repositories.OrganizationRepository.UpdateOrganizationIgnoreEmptyInputParams(ctx, organizationId, eventData)
	} else {
		err = h.repositories.OrganizationRepository.UpdateOrganization(ctx, organizationId, eventData)
		// set customer os id
		customerOsErr := h.setCustomerOsId(ctx, eventData.Tenant, organizationId)
		if customerOsErr != nil {
			tracing.TraceErr(span, customerOsErr)
			h.log.Errorf("Failed to set customer os id for tenant %s organization %s", eventData.Tenant, organizationId)
		}
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
			//var err error
			if eventData.ExternalSystem.Available() {
				innerErr := h.repositories.ExternalSystemRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, organizationId, constants.NodeLabel_Organization, eventData.ExternalSystem)
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
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.PhoneNumberRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.PhoneNumberId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *OrganizationEventHandler) OnEmailLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnEmailLinkedToOrganization")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.EmailRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.EmailId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *OrganizationEventHandler) OnLocationLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnLocationLinkedToOrganization")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationLinkLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.LocationRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.LocationId, eventData.UpdatedAt)

	return err
}

func (h *OrganizationEventHandler) OnDomainLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnDomainLinkedToOrganization")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

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

	err := h.repositories.OrganizationRepository.LinkWithDomain(ctx, eventData.Tenant, organizationId, strings.TrimSpace(eventData.Domain))

	return err
}

func (h *OrganizationEventHandler) OnSocialAddedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnSocialAddedToOrganization")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationAddSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.SocialRepository.MergeSocialFor(ctx, eventData.Tenant, organizationId, "Organization", eventData)

	return err
}

func (h *OrganizationEventHandler) OnOrganizationHide(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationHide")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.HideOrganizationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.OrganizationRepository.SetVisibility(ctx, eventData.Tenant, organizationId, true)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (h *OrganizationEventHandler) OnOrganizationShow(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationShow")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.HideOrganizationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.OrganizationRepository.SetVisibility(ctx, eventData.Tenant, organizationId, false)
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
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRefreshLastTouchpointEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	//fetch the real touchpoint
	//if it doesn't exist, check for the Created Action
	var lastTouchpointId string
	var lastTouchpointAt *time.Time
	var timelineEventNode *dbtype.Node
	var err error

	lastTouchpointAt, lastTouchpointId, err = h.repositories.TimelineEventRepository.CalculateAndGetLastTouchpoint(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to calculate last touchpoint: %v", err.Error())
		span.LogFields(log.Bool("last touchpoint failed", true))
		return nil
	}

	if lastTouchpointAt == nil {
		timelineEventNode, err = h.repositories.ActionRepository.GetSingleAction(ctx, eventData.Tenant, organizationId, entity.ORGANIZATION, entity.ActionCreated)
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
		timelineEventNode, err = h.repositories.TimelineEventRepository.GetTimelineEvent(ctx, eventData.Tenant, lastTouchpointId)
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
	case entity.NodeLabel_PageView:
		timelineEventType = "PAGE_VIEW"
	case entity.NodeLabel_InteractionSession:
		timelineEventType = "INTERACTION_SESSION"
	case entity.NodeLabel_Note:
		timelineEventType = "NOTE"
	case entity.NodeLabel_InteractionEvent:
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
	case entity.NodeLabel_Analysis:
		timelineEventType = "ANALYSIS"
	case entity.NodeLabel_Meeting:
		timelineEventType = "MEETING"
	case entity.NodeLabel_Action:
		timelineEventAction := timelineEvent.(*entity.ActionEntity)
		if timelineEventAction.Type == entity.ActionCreated {
			timelineEventType = "ACTION_CREATED"
		} else {
			timelineEventType = "ACTION"
		}
	case entity.NodeLabel_LogEntry:
		timelineEventType = "LOG_ENTRY"
	case entity.NodeLabel_Issue:
		timelineEventIssue := timelineEvent.(*entity.IssueEntity)
		if timelineEventIssue.CreatedAt.Equal(timelineEventIssue.UpdatedAt) {
			timelineEventType = "ISSUE_CREATED"
		} else {
			timelineEventType = "ISSUE_UPDATED"
		}
	default:
		h.log.Infof("Last touchpoint not available for organization: %s", organizationId)
	}

	if err = h.repositories.OrganizationRepository.UpdateLastTouchpoint(ctx, eventData.Tenant, organizationId, *lastTouchpointAt, lastTouchpointId, timelineEventType); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update last touchpoint for tenant %s, organization %s: %s", eventData.Tenant, organizationId, err.Error())
	}

	return nil
}

func (h *OrganizationEventHandler) OnRefreshArr(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRefreshArr")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRefreshArrEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	if err := h.repositories.OrganizationRepository.UpdateArr(ctx, eventData.Tenant, organizationId); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update arr for tenant %s, organization %s: %s", eventData.Tenant, organizationId, err.Error())
	}

	return nil
}

func (h *OrganizationEventHandler) OnRefreshRenewalSummary(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRefreshRenewalSummary")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRefreshArrEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	openRenewalOpportunityDbNodes, err := h.repositories.OpportunityRepository.GetOpenRenewalOpportunitiesForOrganization(ctx, eventData.Tenant, organizationId)
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

	if err := h.repositories.OrganizationRepository.UpdateRenewalSummary(ctx, eventData.Tenant, organizationId, lowestRenewalLikelihood, renewalLikelihoodOrderPtr, nextRenewalDate); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update arr for tenant %s, organization %s: %s", eventData.Tenant, organizationId, err.Error())
	}

	return nil
}

func getOrderForRenewalLikelihood(likelihood string) int64 {
	switch likelihood {
	case string(opportunitymodel.RenewalLikelihoodStringHigh):
		return constants.RenewalLikelihood_Order_High
	case string(opportunitymodel.RenewalLikelihoodStringMedium):
		return constants.RenewalLikelihood_Order_Medium
	case string(opportunitymodel.RenewalLikelihoodStringLow):
		return constants.RenewalLikelihood_Order_Low
	case string(opportunitymodel.RenewalLikelihoodStringZero):
		return constants.RenewalLikelihood_Order_Zero
	default:
		return 0
	}
}

func (h *OrganizationEventHandler) OnUpsertCustomField(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnUpsertCustomField")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationUpsertCustomField
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	customFieldExists, err := h.repositories.CustomFieldRepository.ExistsById(ctx, eventData.Tenant, eventData.CustomFieldId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to check if custom field exists: %s", err.Error())
		return err
	}
	if !customFieldExists {
		err = h.repositories.CustomFieldRepository.AddCustomFieldToOrganization(ctx, eventData.Tenant, organizationId, eventData)
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
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationAddParentEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	return h.repositories.OrganizationRepository.LinkWithParentOrganization(ctx, eventData.Tenant, organizationId, eventData.ParentOrganizationId, eventData.Type)
}

func (h *OrganizationEventHandler) OnUnlinkFromParentOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnUnlinkFromParentOrganization")
	defer span.Finish()
	setCommonSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRemoveParentEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	return h.repositories.OrganizationRepository.UnlinkParentOrganization(ctx, eventData.Tenant, organizationId, eventData.ParentOrganizationId)
}
