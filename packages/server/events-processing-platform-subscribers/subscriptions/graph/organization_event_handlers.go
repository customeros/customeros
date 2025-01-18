package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	eventsSrv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/organization/aggregate"
	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/organization/events"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/caches"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/constants"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/helper"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/logger"
)

type OrganizationEventHandler struct {
	log             logger.Logger
	grpcClients     *grpc_client.Clients
	cache           caches.Cache
	events          *eventsSrv.EventsService
	neo4j           *neo4jrepository.Repositories
	postgres        *repository.Repositories
	currencyService interfaces.CurrencyService
}

func NewOrganizationEventHandler(
	log logger.Logger,
	grpcClients *grpc_client.Clients,
	cache caches.Cache,
	events *eventsSrv.EventsService,
	neo4j *neo4jrepository.Repositories,
	postgres *repository.Repositories,
	fx interfaces.CurrencyService,
) *OrganizationEventHandler {
	return &OrganizationEventHandler{
		log:             log,
		grpcClients:     grpcClients,
		cache:           cache,
		events:          events,
		neo4j:           neo4j,
		postgres:        postgres,
		currencyService: fx,
	}
}

func (h *OrganizationEventHandler) setCustomerOsId(ctx context.Context, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.setCustomerOsId")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("OrganizationId", organizationId))

	orgDbNode, err := h.neo4j.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)
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
		innerErr := h.postgres.CustomerOsIdsRepository.Reserve(ctx, customerOsIdsEntity)
		if innerErr == nil {
			break
		}
	}
	return h.neo4j.OrganizationWriteRepository.SetCustomerOsIdIfMissing(ctx, tenant, organizationId, customerOsId)
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
	err := h.neo4j.PhoneNumberWriteRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.PhoneNumberId, eventData.Label, eventData.Primary)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	h.events.Publisher.PublishEventCompleted(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

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

	if err := h.neo4j.OrganizationWriteRepository.UpdateArr(ctx, eventData.Tenant, organizationId); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update arr for tenant %s, organization %s: %s", eventData.Tenant, organizationId, err.Error())
	}

	h.events.Publisher.PublishEventCompleted(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

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

	openRenewalOpportunityDbNodes, err := h.neo4j.OpportunityReadRepository.GetActiveRenewalOpportunitiesForOrganization(ctx, eventData.Tenant, organizationId, false)
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
					lowestRenewalLikelihood = utils.ToPtr(opportunity.RenewalDetails.RenewalLikelihood.ToV2().String())
				}
			}
		}
	}

	renewalLikelihoodOrderPtr := utils.ToPtr[int64](renewalLikelihoodOrder)
	if renewalLikelihoodOrder == 0 {
		renewalLikelihoodOrderPtr = nil
	}

	if err := h.neo4j.OrganizationWriteRepository.UpdateRenewalSummary(ctx, eventData.Tenant, organizationId, lowestRenewalLikelihood, renewalLikelihoodOrderPtr, nextRenewalDate); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update arr for tenant %s, organization %s: %s", eventData.Tenant, organizationId, err.Error())
	}

	h.events.Publisher.PublishEventCompleted(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

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

	customFieldExists, err := h.neo4j.CommonReadRepository.ExistsById(ctx, eventData.Tenant, eventData.CustomFieldId, commonmodel.NodeLabelCustomField)
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
		err = h.neo4j.CustomFieldWriteRepository.AddCustomFieldToOrganization(ctx, eventData.Tenant, organizationId, data)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to add custom field to organization: %s", err.Error())
			return err
		}
	} else {
		// TODO implement update custom field
	}
	h.events.Publisher.PublishEventCompleted(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

type ActionOnboardingStatusMetadata struct {
	Status     string `json:"status"`
	Comments   string `json:"comments"`
	UserId     string `json:"userId"`
	ContractId string `json:"contractId"`
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
	err := h.neo4j.BillingProfileWriteRepository.Create(ctx, eventData.Tenant, eventData.BillingProfileId, data)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	h.events.Publisher.PublishEventCompleted(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
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
	err := h.neo4j.BillingProfileWriteRepository.Update(ctx, eventData.Tenant, eventData.BillingProfileId, data)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	h.events.Publisher.PublishEventCompleted(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
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

	err := h.neo4j.BillingProfileWriteRepository.LinkEmailToBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.EmailId, eventData.Primary)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	h.events.Publisher.PublishEventCompleted(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

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

	err := h.neo4j.BillingProfileWriteRepository.UnlinkEmailFromBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.EmailId)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	h.events.Publisher.PublishEventCompleted(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

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

	err := h.neo4j.BillingProfileWriteRepository.LinkLocationToBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.LocationId)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	h.events.Publisher.PublishEventCompleted(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

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

	err := h.neo4j.BillingProfileWriteRepository.UnlinkLocationFromBillingProfile(ctx, eventData.Tenant, organizationId, eventData.BillingProfileId, eventData.LocationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to unlink location %s from billing profile %s: %s", eventData.LocationId, eventData.BillingProfileId, err.Error())
	}

	h.events.Publisher.PublishEventCompleted(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

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

	organizationDbNode, err := h.neo4j.OrganizationReadRepository.GetOrganization(ctx, eventData.Tenant, organizationId)
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

	h.events.Publisher.PublishEventCompleted(ctx, eventData.Tenant, organizationId, commonmodel.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (h *OrganizationEventHandler) deriveChurnedDate(ctx context.Context, tenant string, organizationEntity *neo4jentity.OrganizationEntity, span opentracing.Span) error {
	// get all contracts for organization
	orgContracts, err := h.neo4j.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{organizationEntity.ID})
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
		err = h.neo4j.OrganizationWriteRepository.UpdateTimeProperty(ctx, tenant, organizationEntity.ID, "derivedChurnedAt", endedAt)
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
	orgContracts, err := h.neo4j.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{organizationEntity.ID})
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
		tenantSettingsDbNode, err := h.neo4j.TenantReadRepository.GetTenantSettings(ctx, tenant)
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
			rate, err := h.currencyService.GetRate(ctx, contract.Currency.String(), ltvCurrency)
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
	err = h.neo4j.OrganizationWriteRepository.UpdateFloatProperty(ctx, tenant, organizationEntity.ID, "derivedLtv", truncatedLtv)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update ltv for organization %s: %s", organizationEntity.ID, err.Error())
	}

	// set ltv currency
	if ltvCurrency != "" {
		err = h.neo4j.OrganizationWriteRepository.UpdateStringProperty(ctx, tenant, organizationEntity.ID, "derivedLtvCurrency", ltvCurrency)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to update ltv currency for organization %s: %s", organizationEntity.ID, err.Error())
		}
	}

	span.LogFields(log.String("result.ltv", fmt.Sprintf("%f", truncatedLtv)))
	span.LogFields(log.String("result.ltvCurrency", ltvCurrency))

	return nil
}
