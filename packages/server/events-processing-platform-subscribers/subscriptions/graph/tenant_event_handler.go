package graph

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"

	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	tenant "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type TenantEventHandler struct {
	log      logger.Logger
	services *service.Services
}

func NewTenantEventHandler(log logger.Logger, services *service.Services) *TenantEventHandler {
	return &TenantEventHandler{
		log:      log,
		services: services,
	}
}

func (h *TenantEventHandler) OnAddBillingProfileV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantEventHandler.OnAddBillingProfileV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.TenantBillingProfileCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	tenantName := tenant.GetTenantName(evt.GetAggregateID())
	span.SetTag(tracing.SpanTagEntityId, tenantName)

	data := neo4jrepository.TenantBillingProfileCreateFields{
		Id:        eventData.Id,
		CreatedAt: eventData.CreatedAt,
		SourceFields: neo4jmodel.SourceFields{
			Source:    helper.GetSource(eventData.SourceFields.Source),
			AppSource: helper.GetAppSource(eventData.SourceFields.AppSource),
		},
		Phone:                  eventData.Phone,
		LegalName:              eventData.LegalName,
		AddressLine1:           eventData.AddressLine1,
		AddressLine2:           eventData.AddressLine2,
		AddressLine3:           eventData.AddressLine3,
		Locality:               eventData.Locality,
		Country:                eventData.Country,
		Region:                 eventData.Region,
		Zip:                    eventData.Zip,
		VatNumber:              eventData.VatNumber,
		SendInvoicesFrom:       eventData.SendInvoicesFrom,
		SendInvoicesBcc:        eventData.SendInvoicesBcc,
		CanPayWithPigeon:       eventData.CanPayWithPigeon,
		CanPayWithBankTransfer: eventData.CanPayWithBankTransfer,
		Check:                  eventData.Check,
	}
	err := h.services.CommonServices.Neo4jRepositories.TenantWriteRepository.CreateTenantBillingProfile(ctx, tenantName, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return err
}

func (h *TenantEventHandler) OnUpdateBillingProfileV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantEventHandler.OnAddBillingProfileV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.TenantBillingProfileUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	tenantName := tenant.GetTenantName(evt.GetAggregateID())
	span.SetTag(tracing.SpanTagEntityId, tenantName)

	data := neo4jrepository.TenantBillingProfileUpdateFields{
		Id:                           eventData.Id,
		Phone:                        eventData.Phone,
		LegalName:                    eventData.LegalName,
		AddressLine1:                 eventData.AddressLine1,
		AddressLine2:                 eventData.AddressLine2,
		AddressLine3:                 eventData.AddressLine3,
		Locality:                     eventData.Locality,
		Country:                      eventData.Country,
		Region:                       eventData.Region,
		Zip:                          eventData.Zip,
		VatNumber:                    eventData.VatNumber,
		SendInvoicesFrom:             eventData.SendInvoicesFrom,
		SendInvoicesBcc:              eventData.SendInvoicesBcc,
		CanPayWithPigeon:             eventData.CanPayWithPigeon,
		CanPayWithBankTransfer:       eventData.CanPayWithBankTransfer,
		Check:                        eventData.Check,
		UpdatePhone:                  eventData.UpdatePhone(),
		UpdateAddressLine1:           eventData.UpdateAddressLine1(),
		UpdateAddressLine2:           eventData.UpdateAddressLine2(),
		UpdateAddressLine3:           eventData.UpdateAddressLine3(),
		UpdateLocality:               eventData.UpdateLocality(),
		UpdateCountry:                eventData.UpdateCountry(),
		UpdateRegion:                 eventData.UpdateRegion(),
		UpdateZip:                    eventData.UpdateZip(),
		UpdateLegalName:              eventData.UpdateLegalName(),
		UpdateVatNumber:              eventData.UpdateVatNumber(),
		UpdateSendInvoicesFrom:       eventData.UpdateSendInvoicesFrom(),
		UpdateSendInvoicesBcc:        eventData.UpdateSendInvoicesBcc(),
		UpdateCanPayWithPigeon:       eventData.UpdateCanPayWithPigeon(),
		UpdateCanPayWithBankTransfer: eventData.UpdateCanPayWithBankTransfer(),
		UpdateCheck:                  eventData.UpdateCheck(),
	}
	err := h.services.CommonServices.Neo4jRepositories.TenantWriteRepository.UpdateTenantBillingProfile(ctx, tenantName, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return err
}

func (h *TenantEventHandler) OnUpdateTenantSettingsV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantEventHandler.OnUpdateTenantSettingsV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.TenantSettingsUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	tenantName := tenant.GetTenantName(evt.GetAggregateID())
	span.SetTag(tracing.SpanTagEntityId, tenantName)

	data := neo4jrepository.TenantSettingsFields{
		LogoRepositoryFileId:       eventData.LogoRepositoryFileId,
		InvoicingEnabled:           eventData.InvoicingEnabled,
		InvoicingPostpaid:          eventData.InvoicingPostpaid,
		WorkspaceLogo:              eventData.WorkspaceLogo,
		WorkspaceName:              eventData.WorkspaceName,
		BaseCurrency:               neo4jenum.DecodeCurrency(eventData.BaseCurrency),
		UpdateInvoicingEnabled:     eventData.UpdateInvoicingEnabled(),
		UpdateBaseCurrency:         eventData.UpdateBaseCurrency(),
		UpdateInvoicingPostpaid:    eventData.UpdateInvoicingPostpaid(),
		UpdateLogoRepositoryFileId: eventData.UpdateLogoRepositoryFileId(),
		UpdateWorkspaceLogo:        eventData.UpdateWorkspaceLogo(),
		UpdateWorkspaceName:        eventData.UpdateWorkspaceName(),
	}
	err := h.services.CommonServices.Neo4jRepositories.TenantWriteRepository.UpdateTenantSettings(ctx, tenantName, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return err
}
