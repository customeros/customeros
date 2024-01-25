package graph

import (
	"context"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	tenant "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type TenantEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewTenantEventHandler(log logger.Logger, repositories *repository.Repositories) *TenantEventHandler {
	return &TenantEventHandler{
		log:          log,
		repositories: repositories,
	}
}

func (h *TenantEventHandler) OnAddBillingProfile(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantEventHandler.OnAddBillingProfile")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.CreateTenantBillingProfileEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	tenantName := tenant.GetTenantName(evt.GetAggregateID())
	span.SetTag(tracing.SpanTagEntityId, tenantName)

	data := neo4jrepository.TenantBillingProfileCreateFields{
		Id:        eventData.Id,
		CreatedAt: eventData.CreatedAt,
		SourceFields: neo4jmodel.Source{
			Source:    helper.GetSource(eventData.SourceFields.Source),
			AppSource: helper.GetAppSource(eventData.SourceFields.AppSource),
		},
		Email:                             eventData.Email,
		Phone:                             eventData.Phone,
		LegalName:                         eventData.LegalName,
		AddressLine1:                      eventData.AddressLine1,
		AddressLine2:                      eventData.AddressLine2,
		AddressLine3:                      eventData.AddressLine3,
		Locality:                          eventData.Locality,
		Country:                           eventData.Country,
		Zip:                               eventData.Zip,
		DomesticPaymentsBankInfo:          eventData.DomesticPaymentsBankInfo,
		DomesticPaymentsBankName:          eventData.DomesticPaymentsBankName,
		DomesticPaymentsAccountNumber:     eventData.DomesticPaymentsAccountNumber,
		DomesticPaymentsSortCode:          eventData.DomesticPaymentsSortCode,
		InternationalPaymentsBankInfo:     eventData.InternationalPaymentsBankInfo,
		InternationalPaymentsSwiftBic:     eventData.InternationalPaymentsSwiftBic,
		InternationalPaymentsBankName:     eventData.InternationalPaymentsBankName,
		InternationalPaymentsBankAddress:  eventData.InternationalPaymentsBankAddress,
		InternationalPaymentsInstructions: eventData.InternationalPaymentsInstructions,
	}
	err := h.repositories.Neo4jRepositories.TenantWriteRepository.CreateTenantBillingProfile(ctx, tenantName, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return err
}
