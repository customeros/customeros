package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type OrganizationEventHandler struct {
	log         logger.Logger
	grpcClients *grpc_client.Clients
	cache       caches.Cache
	services    *service.Services
}

func NewOrganizationEventHandler(log logger.Logger, services *service.Services, grpcClients *grpc_client.Clients, cache caches.Cache) *OrganizationEventHandler {
	return &OrganizationEventHandler{
		log:         log,
		grpcClients: grpcClients,
		cache:       cache,
		services:    services,
	}
}

type eventMetadata struct {
	UserId string `json:"user-id"`
}

// DEPRECATED
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

	session := utils.NewNeo4jWriteSession(ctx, *h.services.CommonServices.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var err error
		data := neo4jrepository.OrganizationCreateFields{
			AggregateVersion: evt.Version,
			CreatedAt:        eventData.CreatedAt,
			SourceFields: neo4jmodel.SourceFields{
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
			Employees:          eventData.Employees,
			Market:             eventData.Market,
			LastFundingRound:   eventData.LastFundingRound,
			LastFundingAmount:  eventData.LastFundingAmount,
			ReferenceId:        eventData.ReferenceId,
			Note:               eventData.Note,
			LogoUrl:            eventData.LogoUrl,
			IconUrl:            eventData.IconUrl,
			Headquarters:       eventData.Headquarters,
			YearFounded:        eventData.YearFounded,
			EmployeeGrowthRate: eventData.EmployeeGrowthRate,
			SlackChannelId:     eventData.SlackChannelId,
			Relationship:       neo4jenum.DecodeOrganizationRelationship(eventData.Relationship),
			Stage:              neo4jenum.DecodeOrganizationStage(eventData.Stage),
			LeadSource:         eventData.LeadSource,
			IcpFit:             eventData.IcpFit,
		}
		err = h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.CreateOrganizationInTx(ctx, tx, eventData.Tenant, organizationId, data)
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
			err = h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, organizationId, commonmodel.NodeLabelOrganization, externalSystemData)
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
			err = h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.ReplaceOwner(ctx, nil, eventData.Tenant, organizationId, evtMetadata.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to replace owner of organization %s with user %s", organizationId, evtMetadata.UserId)
			}
		}
	}

	// Set create action
	_, err = h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.MergeByActionType(ctx, nil, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, neo4jenum.ActionCreated, "", "", eventData.CreatedAt, constants.AppSourceEventProcessingPlatformSubscribers)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed creating likelihood update action for organization %s: %s", organizationId, err.Error())
	}

	// set domain
	if eventData.Website != "" {
		primaryDomain, _ := h.services.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, eventData.Website)
		if primaryDomain != "" {
			err = h.services.CommonServices.OrganizationService.LinkWithDomain(ctx, nil, organizationId, primaryDomain)
		}
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to add domain from website %s: %s", eventData.Website, err.Error())
		}
	}

	// Request last touch point update
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
		return h.grpcClients.OrganizationClient.RefreshLastTouchpoint(ctx, &organizationpb.OrganizationIdGrpcRequest{
			Tenant:         eventData.Tenant,
			OrganizationId: organizationId,
			AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while refreshing last touchpoint for organization %s: %s", organizationId, err.Error())
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithCreate())

	return nil
}

func (h *OrganizationEventHandler) setCustomerOsId(ctx context.Context, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.setCustomerOsId")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("OrganizationId", organizationId))

	orgDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)
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
		innerErr := h.services.CommonServices.PostgresRepositories.CustomerOsIdsRepository.Reserve(ctx, customerOsIdsEntity)
		if innerErr == nil {
			break
		}
	}
	return h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.SetCustomerOsIdIfMissing(ctx, tenant, organizationId, customerOsId)
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
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	var beforeOrganizationEntity, afterOrganizationEntity neo4jentity.OrganizationEntity
	existingOrganization, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if existingOrganization != nil {
		beforeOrganizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(existingOrganization)
	}

	data := neo4jrepository.OrganizationUpdateFields{
		AggregateVersion:         evt.Version,
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
		Employees:                eventData.Employees,
		Market:                   eventData.Market,
		LastFundingRound:         eventData.LastFundingRound,
		LastFundingAmount:        eventData.LastFundingAmount,
		ReferenceId:              eventData.ReferenceId,
		Note:                     eventData.Note,
		LogoUrl:                  eventData.LogoUrl,
		IconUrl:                  eventData.IconUrl,
		Headquarters:             eventData.Headquarters,
		YearFounded:              eventData.YearFounded,
		EmployeeGrowthRate:       eventData.EmployeeGrowthRate,
		SlackChannelId:           eventData.SlackChannelId,
		EnrichDomain:             eventData.EnrichDomain,
		EnrichSource:             eventData.EnrichSource,
		Source:                   helper.GetSource(eventData.Source),
		Relationship:             neo4jenum.DecodeOrganizationRelationship(eventData.Relationship),
		Stage:                    neo4jenum.DecodeOrganizationStage(eventData.Stage),
		IcpFit:                   eventData.IcpFit,
		UpdateName:               eventData.UpdateName(),
		UpdateDescription:        eventData.UpdateDescription(),
		UpdateHide:               eventData.UpdateHide(),
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
		UpdateIconUrl:            eventData.UpdateIconUrl(),
		UpdateEmployeeGrowthRate: eventData.UpdateEmployeeGrowthRate(),
		UpdateSlackChannelId:     eventData.UpdateSlackChannelId(),
		UpdateRelationship:       eventData.UpdateRelationship(),
		UpdateStage:              eventData.UpdateStage(),
		UpdateIcpFit:             eventData.UpdateIcpFit(),
	}

	err = h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.UpdateOrganization(ctx, eventData.Tenant, organizationId, data)

	// set customer os id if not set
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

	// link with external system
	if eventData.ExternalSystem.Available() {
		session := utils.NewNeo4jWriteSession(ctx, *h.services.CommonServices.Neo4jRepositories.Neo4jDriver)
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
				innerErr := h.services.CommonServices.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, organizationId, commonmodel.NodeLabelOrganization, externalSystemData)
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

	// set slack channel id for unthreaded issues
	if eventData.UpdateSlackChannelId() {
		if beforeOrganizationEntity.ID != "" && beforeOrganizationEntity.SlackChannelId != eventData.SlackChannelId {
			if eventData.SlackChannelId == "" {
				err := h.services.CommonServices.Neo4jRepositories.IssueWriteRepository.RemoveReportedByOrganizationWithGroupId(ctx, eventData.Tenant, organizationId, beforeOrganizationEntity.SlackChannelId)
				if err != nil {
					tracing.TraceErr(span, err)
					h.log.Errorf("Failed to remove reported by organization with groupId %s: %s", beforeOrganizationEntity.SlackChannelId, err.Error())
				}
			} else {
				err := h.services.CommonServices.Neo4jRepositories.IssueWriteRepository.ReportedByOrganizationWithGroupId(ctx, eventData.Tenant, organizationId, eventData.SlackChannelId)
				if err != nil {
					tracing.TraceErr(span, err)
					h.log.Errorf("Failed to mark reported by organization with groupId %s: %s", eventData.SlackChannelId, err.Error())
				}
			}
		}
	}

	updatedOrganization, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	afterOrganizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(updatedOrganization)

	if beforeOrganizationEntity.Website != afterOrganizationEntity.Website {
		primaryDomain, _ := h.services.CommonServices.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, eventData.Website)
		if primaryDomain != "" {
			err = h.services.CommonServices.OrganizationService.LinkWithDomain(ctx, nil, organizationId, primaryDomain)
		}
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to add domain from website %s: %s", eventData.Website, err.Error())
		}
	}

	if beforeOrganizationEntity.Stage != afterOrganizationEntity.Stage {
		h.handleStageChange(ctx, eventData.Tenant, &beforeOrganizationEntity, &afterOrganizationEntity)
	}

	if eventData.AppSource != constants.AppSourceCustomerOsApi {
		utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())
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
	err := h.services.CommonServices.Neo4jRepositories.PhoneNumberWriteRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.PhoneNumberId, eventData.Label, eventData.Primary)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
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
	err := h.services.CommonServices.Neo4jRepositories.LocationWriteRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.LocationId)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return err
}

func (h *OrganizationEventHandler) OnDomainUnlinkedFromOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnDomainUnlinkedFromOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationUnlinkDomainEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	err := h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.UnlinkFromDomain(ctx, eventData.Tenant, organizationId, eventData.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
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
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	data := neo4jrepository.SocialFields{
		SocialId:       eventData.SocialId,
		Url:            eventData.Url,
		Alias:          eventData.Alias,
		ExternalId:     eventData.ExternalId,
		FollowersCount: eventData.FollowersCount,
		CreatedAt:      eventData.CreatedAt,
		SourceFields: neo4jmodel.SourceFields{
			Source:        helper.GetSource(eventData.Source),
			SourceOfTruth: helper.GetSource(eventData.SourceOfTruth),
			AppSource:     helper.GetSource(eventData.AppSource),
		},
	}
	err := h.services.CommonServices.Neo4jRepositories.SocialWriteRepository.MergeSocialForEntity(ctx, eventData.Tenant, organizationId, commonmodel.NodeLabelOrganization, data)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *OrganizationEventHandler) OnSocialRemovedFromOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnSocialRemovedFromOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRemoveSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	if eventData.SocialId != "" {
		err := h.services.CommonServices.Neo4jRepositories.SocialWriteRepository.RemoveSocialForEntityById(ctx, eventData.Tenant, organizationId, commonmodel.NodeLabelOrganization, eventData.SocialId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil
		}
	} else {
		err := h.services.CommonServices.Neo4jRepositories.SocialWriteRepository.RemoveSocialForEntityByUrl(ctx, eventData.Tenant, organizationId, commonmodel.NodeLabelOrganization, eventData.Url)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil
		}
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *OrganizationEventHandler) OnOrganizationShow(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationShow")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.ShowOrganizationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.SetVisibility(ctx, eventData.Tenant, organizationId, false)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	// set customer os id
	customerOsErr := h.setCustomerOsId(ctx, eventData.Tenant, organizationId)
	if customerOsErr != nil {
		tracing.TraceErr(span, customerOsErr)
		h.log.Errorf("Failed to set customer os id for tenant %s organization %s", eventData.Tenant, organizationId)
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return err
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

	if err := h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.UpdateArr(ctx, eventData.Tenant, organizationId); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update arr for tenant %s, organization %s: %s", eventData.Tenant, organizationId, err.Error())
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *OrganizationEventHandler) OnRefreshRenewalSummaryV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRefreshRenewalSummaryV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRefreshArrEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	openRenewalOpportunityDbNodes, err := h.services.CommonServices.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunitiesForOrganization(ctx, eventData.Tenant, organizationId, false)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to get open renewal opportunities for organization %s: %s", organizationId, err.Error())
		return nil
	}
	var nextRenewalDate *time.Time
	var lowestRenewalLikelihood *string
	var renewalLikelihoodOrder int64
	if len(openRenewalOpportunityDbNodes) > 0 {
		opportunities := make([]neo4jentity.OpportunityEntity, len(openRenewalOpportunityDbNodes))
		for _, opportunityDbNode := range openRenewalOpportunityDbNodes {
			opportunities = append(opportunities, *neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode))
		}
		for _, opportunity := range opportunities {
			if opportunity.RenewalDetails.RenewedAt != nil && opportunity.RenewalDetails.RenewedAt.After(utils.Now()) {
				if nextRenewalDate == nil || opportunity.RenewalDetails.RenewedAt.Before(*nextRenewalDate) {
					nextRenewalDate = opportunity.RenewalDetails.RenewedAt
				}
			}
			if opportunity.RenewalDetails.RenewalLikelihood != "" {
				order := getOrderForRenewalLikelihood(opportunity.RenewalDetails.RenewalLikelihood.String())
				if renewalLikelihoodOrder == 0 || renewalLikelihoodOrder > order {
					renewalLikelihoodOrder = order
					lowestRenewalLikelihood = utils.ToPtr(opportunity.RenewalDetails.RenewalLikelihood.String())
				}
			}
		}
	}

	renewalLikelihoodOrderPtr := utils.ToPtr[int64](renewalLikelihoodOrder)
	if renewalLikelihoodOrder == 0 {
		renewalLikelihoodOrderPtr = nil
	}

	if err := h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.UpdateRenewalSummary(ctx, eventData.Tenant, organizationId, lowestRenewalLikelihood, renewalLikelihoodOrderPtr, nextRenewalDate); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update arr for tenant %s, organization %s: %s", eventData.Tenant, organizationId, err.Error())
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

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

	customFieldExists, err := h.services.CommonServices.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, eventData.Tenant, eventData.CustomFieldId, commonmodel.NodeLabelCustomField)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to check if custom field exists: %s", err.Error())
		return err
	}
	if !customFieldExists {
		data := neo4jrepository.CustomFieldCreateFields{
			CreatedAt:           eventData.CreatedAt,
			ExistsInEventStore:  eventData.ExistsInEventStore,
			TemplateId:          eventData.TemplateId,
			CustomFieldId:       eventData.CustomFieldId,
			CustomFieldName:     eventData.CustomFieldName,
			CustomFieldDataType: eventData.CustomFieldDataType,
			CustomFieldValue:    eventData.CustomFieldValue,
			SourceFields: neo4jmodel.SourceFields{
				Source:        helper.GetSource(eventData.Source),
				SourceOfTruth: helper.GetSource(eventData.SourceOfTruth),
				AppSource:     helper.GetSource(eventData.AppSource),
			},
		}
		err = h.services.CommonServices.Neo4jRepositories.CustomFieldWriteRepository.AddCustomFieldToOrganization(ctx, eventData.Tenant, organizationId, data)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to add custom field to organization: %s", err.Error())
			return err
		}
	} else {
		//TODO implement update custom field
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

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
	err := h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.LinkWithParentOrganization(ctx, eventData.Tenant, organizationId, eventData.ParentOrganizationId, eventData.Type)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())
	return nil
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
	err := h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.UnlinkParentOrganization(ctx, eventData.Tenant, organizationId, eventData.ParentOrganizationId)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())
	return nil
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

	organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, eventData.Tenant, organizationId)
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

	err = h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.UpdateOnboardingStatus(ctx, eventData.Tenant, organizationId, eventData.Status, eventData.Comments, getOrderForOnboardingStatus(eventData.Status), eventData.UpdatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update onboarding status for organization %s: %s", organizationId, err.Error())
		return err
	}

	if eventData.CausedByContractId != "" {
		err = h.services.CommonServices.Neo4jRepositories.ContractWriteRepository.ContractCausedOnboardingStatusChange(ctx, eventData.Tenant, eventData.CausedByContractId)
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

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

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
		userDbNode, err := h.services.CommonServices.Neo4jRepositories.UserReadRepository.GetUserById(ctx, eventData.Tenant, eventData.UpdatedByUserId)
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
	_, err := h.services.CommonServices.Neo4jRepositories.ActionWriteRepository.CreateWithProperties(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, neo4jenum.ActionOnboardingStatusChanged, message, metadata, eventData.UpdatedAt, constants.AppSourceEventProcessingPlatformSubscribers, extraActionProperties)
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
		SourceFields: neo4jmodel.SourceFields{
			Source:    helper.GetSource(eventData.SourceFields.Source),
			AppSource: helper.GetSource(eventData.SourceFields.AppSource),
		},
	}
	err := h.services.CommonServices.Neo4jRepositories.BillingProfileWriteRepository.Create(ctx, eventData.Tenant, eventData.BillingProfileId, data)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())
	return nil
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
		LegalName:       eventData.LegalName,
		TaxId:           eventData.TaxId,
		UpdateLegalName: eventData.UpdateLegalName(),
		UpdateTaxId:     eventData.UpdateTaxId(),
	}
	err := h.services.CommonServices.Neo4jRepositories.BillingProfileWriteRepository.Update(ctx, eventData.Tenant, eventData.BillingProfileId, data)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())
	return nil
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

	err := h.services.CommonServices.Neo4jRepositories.BillingProfileWriteRepository.LinkEmailToBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.EmailId, eventData.Primary)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
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

	err := h.services.CommonServices.Neo4jRepositories.BillingProfileWriteRepository.UnlinkEmailFromBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.EmailId)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
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

	err := h.services.CommonServices.Neo4jRepositories.BillingProfileWriteRepository.LinkLocationToBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.LocationId)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
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

	err := h.services.CommonServices.Neo4jRepositories.BillingProfileWriteRepository.UnlinkLocationFromBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.LocationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to unlink location %s from billing profile %s: %s", eventData.LocationId, eventData.BillingProfileId, err.Error())
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *OrganizationEventHandler) OnRefreshDerivedDataV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRefreshDerivedDataV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationRefreshDerivedData
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	tenant := eventData.Tenant
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.SetTag(tracing.SpanTagTenant, tenant)

	organizationDbNode, err := h.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to get organization %s: %s", organizationId, err.Error())
		return err
	}
	organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	err = h.deriveChurnedDate(ctx, tenant, organizationEntity, span)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	err = h.deriveLtv(ctx, tenant, organizationEntity, span)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *OrganizationEventHandler) deriveChurnedDate(ctx context.Context, tenant string, organizationEntity *neo4jentity.OrganizationEntity, span opentracing.Span) error {
	// get all contracts for organization
	orgContracts, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{organizationEntity.ID})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting contracts for organization %s: %s", organizationEntity.ID, err.Error())
		return err
	}

	orgContractEntities := []neo4jentity.ContractEntity{}
	for _, orgContract := range orgContracts {
		orgContractEntities = append(orgContractEntities, *neo4jmapper.MapDbNodeToContractEntity(orgContract.Node))
	}
	endedContractFound := false
	nonEndedContractFound := false
	var endedAt *time.Time

	for _, contract := range orgContractEntities {
		if contract.ContractStatus == neo4jenum.ContractStatusDraft {
			continue
		}
		if contract.ContractStatus == neo4jenum.ContractStatusEnded {
			endedContractFound = true
			if contract.EndedAt != nil && (endedAt == nil || contract.EndedAt.After(*endedAt)) {
				endedAt = contract.EndedAt
			}
		}
		if contract.ContractStatus != neo4jenum.ContractStatusEnded && contract.ContractStatus != neo4jenum.ContractStatusDraft {
			nonEndedContractFound = true
			break
		}
	}

	if nonEndedContractFound {
		span.LogFields(log.String("result", "no non-ended contracts found"))
		return nil
	}

	if endedContractFound && endedAt != nil {
		err = h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.UpdateTimeProperty(ctx, tenant, organizationEntity.ID, "derivedChurnedAt", endedAt)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to update churn date for organization %s: %s", organizationEntity.ID, err.Error())
			return err
		}
	}

	return nil
}

func (h *OrganizationEventHandler) deriveLtv(ctx context.Context, tenant string, organizationEntity *neo4jentity.OrganizationEntity, span opentracing.Span) error {
	// get all contracts for organization
	orgContracts, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{organizationEntity.ID})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting contracts for organization %s: %s", organizationEntity.ID, err.Error())
		return err
	}

	var orgContractEntities []neo4jentity.ContractEntity
	for _, orgContract := range orgContracts {
		orgContractEntities = append(orgContractEntities, *neo4jmapper.MapDbNodeToContractEntity(orgContract.Node))
	}

	// check multiple currencies
	currencySet := make(map[string]struct{})
	for _, contract := range orgContractEntities {
		if contract.Ltv != 0 && contract.Currency.String() != "" {
			currencySet[contract.Currency.String()] = struct{}{}
		}
	}
	multipleCurrencies := len(currencySet) > 1

	ltvCurrency := ""
	if multipleCurrencies {
		// get tenant base currency
		tenantSettingsDbNode, err := h.services.CommonServices.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to get tenant settings for tenant %s: %s", tenant, err.Error())
		}
		tenantSettings := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettingsDbNode)
		ltvCurrency = tenantSettings.BaseCurrency.String()
	} else if len(currencySet) == 1 {
		ltvCurrency = orgContractEntities[0].Currency.String()
	}

	ltv := 0.0
	for _, contract := range orgContractEntities {
		if ltvCurrency == "" || contract.Currency.String() == ltvCurrency {
			ltv += contract.Ltv
		} else {
			rate, err := h.services.CommonServices.CurrencyService.GetRate(ctx, contract.Currency.String(), ltvCurrency)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to get rate for currency %s: %s", contract.Currency.String(), err.Error())
				continue
			}
			ltv += contract.Ltv * rate
		}
	}

	// set ltv
	truncatedLtv := utils.TruncateFloat64(ltv, 2)
	err = h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.UpdateFloatProperty(ctx, tenant, organizationEntity.ID, "derivedLtv", truncatedLtv)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update ltv for organization %s: %s", organizationEntity.ID, err.Error())
	}

	// set ltv currency
	if ltvCurrency != "" {
		err = h.services.CommonServices.Neo4jRepositories.OrganizationWriteRepository.UpdateStringProperty(ctx, tenant, organizationEntity.ID, "derivedLtvCurrency", ltvCurrency)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to update ltv currency for organization %s: %s", organizationEntity.ID, err.Error())
		}
	}

	span.LogFields(log.String("result.ltv", fmt.Sprintf("%f", truncatedLtv)))
	span.LogFields(log.String("result.ltvCurrency", ltvCurrency))

	return nil
}

func (h *OrganizationEventHandler) handleStageChange(ctx context.Context, tenant string, orgBeforeUpdate, orgAfterUpdate *neo4jentity.OrganizationEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.handleStageChange")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, orgBeforeUpdate.ID)

	if orgBeforeUpdate.Stage == orgAfterUpdate.Stage {
		return
	}

	// Create default opportunity for organization
	if orgAfterUpdate.Stage == neo4jenum.Engaged {
		// check if organization has any opportunity
		opportunities, err := h.services.CommonServices.Neo4jRepositories.OpportunityReadRepository.GetForOrganizations(ctx, tenant, []string{orgAfterUpdate.ID})
		if err != nil {
			tracing.TraceErr(span, err)
		}
		if len(opportunities) > 0 {
			span.LogFields(log.String("result", "skip opportunity creation, opportunity already exists"))
			return
		}

		// check if organization has any contract
		contracts, err := h.services.CommonServices.Neo4jRepositories.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{orgAfterUpdate.ID})
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
		if len(contracts) > 0 {
			span.LogFields(log.String("result", "skip opportunity creation, contract already exists"))
			return
		}

		// get organization owner
		ownerDbNode, err := h.services.CommonServices.Neo4jRepositories.UserReadRepository.GetOwnerForOrganization(ctx, tenant, orgAfterUpdate.ID)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		ownerId := ""
		if ownerDbNode != nil {
			owner := neo4jmapper.MapDbNodeToUserEntity(ownerDbNode)
			ownerId = owner.Id
		}

		// get default opportunity stage
		tenantSettingStages, err := h.services.CommonServices.PostgresRepositories.TenantSettingsOpportunityStageRepository.GetOrInitialize(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while getting tenant setting opportunity stages: %s", err.Error())
			return
		}
		// get first visible stage
		defaultStage := tenantSettingStages[0].Value
		for _, stage := range tenantSettingStages {
			if stage.Visible {
				defaultStage = stage.Value
				break
			}
		}

		// create default opportunity
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
			return h.grpcClients.OpportunityClient.CreateOpportunity(ctx, &opportunitypb.CreateOpportunityGrpcRequest{
				Tenant:         tenant,
				OrganizationId: orgAfterUpdate.ID,
				Name:           orgAfterUpdate.Name,
				OwnerUserId:    ownerId,
				InternalType:   opportunitypb.OpportunityInternalType_NBO,
				InternalStage:  opportunitypb.OpportunityInternalStage_OPEN,
				ExternalStage:  defaultStage,
				SourceFields: &commonpb.SourceFields{
					Source:    constants.SourceOpenline,
					AppSource: constants.AppSourceEventProcessingPlatformSubscribers,
				},
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Error while creating default opportunity for organization %s: %s", orgAfterUpdate.ID, err.Error())
		}
	}
}

func (h *OrganizationEventHandler) OnLocationAddedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnLocationAddedToOrganization")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.OrganizationAddLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	data := neo4jrepository.LocationCreateFields{
		RawAddress: eventData.RawAddress,
		Name:       eventData.Name,
		CreatedAt:  eventData.CreatedAt,
		SourceFields: neo4jmodel.SourceFields{
			Source:        helper.GetSource(eventData.Source),
			SourceOfTruth: helper.GetSource(eventData.SourceOfTruth),
			AppSource:     helper.GetSource(eventData.AppSource),
		},
		AddressDetails: neo4jrepository.AddressDetails{
			Latitude:      eventData.Latitude,
			Longitude:     eventData.Longitude,
			Country:       eventData.Country,
			CountryCodeA2: eventData.CountryCodeA2,
			CountryCodeA3: eventData.CountryCodeA3,
			Region:        eventData.Region,
			District:      eventData.District,
			Locality:      eventData.Locality,
			Street:        eventData.Street,
			Address:       eventData.AddressLine1,
			Address2:      eventData.AddressLine2,
			Zip:           eventData.ZipCode,
			AddressType:   eventData.AddressType,
			HouseNumber:   eventData.HouseNumber,
			PostalCode:    eventData.PostalCode,
			PlusFour:      eventData.PlusFour,
			Commercial:    eventData.Commercial,
			Predirection:  eventData.Predirection,
			TimeZone:      eventData.TimeZone,
			UtcOffset:     eventData.UtcOffset,
		},
	}

	err := h.services.CommonServices.Neo4jRepositories.LocationWriteRepository.CreateLocation(ctx, eventData.Tenant, eventData.LocationId, data)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while creating location %s: %s", eventData.LocationId, err.Error())
		return err
	}
	err = h.services.CommonServices.Neo4jRepositories.LocationWriteRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.LocationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("error while linking location %s to organization %s: %s", eventData.LocationId, organizationId, err.Error())
		return err
	}

	utils.EventCompleted(ctx, eventData.Tenant, commonmodel.ORGANIZATION.String(), organizationId, h.grpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}
