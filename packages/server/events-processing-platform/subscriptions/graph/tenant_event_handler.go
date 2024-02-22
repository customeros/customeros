package graph

import (
	"context"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
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
		VatNumber:                         eventData.VatNumber,
		SendInvoicesFrom:                  eventData.SendInvoicesFrom,
		CanPayWithCard:                    eventData.CanPayWithCard,
		CanPayWithDirectDebitSEPA:         eventData.CanPayWithDirectDebitSEPA,
		CanPayWithDirectDebitACH:          eventData.CanPayWithDirectDebitACH,
		CanPayWithDirectDebitBacs:         eventData.CanPayWithDirectDebitBacs,
		CanPayWithPigeon:                  eventData.CanPayWithPigeon,
	}
	err := h.repositories.Neo4jRepositories.TenantWriteRepository.CreateTenantBillingProfile(ctx, tenantName, data)
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
		Id:                                  eventData.Id,
		UpdatedAt:                           eventData.UpdatedAt,
		Email:                               eventData.Email,
		Phone:                               eventData.Phone,
		LegalName:                           eventData.LegalName,
		AddressLine1:                        eventData.AddressLine1,
		AddressLine2:                        eventData.AddressLine2,
		AddressLine3:                        eventData.AddressLine3,
		Locality:                            eventData.Locality,
		Country:                             eventData.Country,
		Zip:                                 eventData.Zip,
		DomesticPaymentsBankInfo:            eventData.DomesticPaymentsBankInfo,
		InternationalPaymentsBankInfo:       eventData.InternationalPaymentsBankInfo,
		VatNumber:                           eventData.VatNumber,
		SendInvoicesFrom:                    eventData.SendInvoicesFrom,
		CanPayWithCard:                      eventData.CanPayWithCard,
		CanPayWithDirectDebitSEPA:           eventData.CanPayWithDirectDebitSEPA,
		CanPayWithDirectDebitACH:            eventData.CanPayWithDirectDebitACH,
		CanPayWithDirectDebitBacs:           eventData.CanPayWithDirectDebitBacs,
		CanPayWithPigeon:                    eventData.CanPayWithPigeon,
		UpdateEmail:                         eventData.UpdateEmail(),
		UpdatePhone:                         eventData.UpdatePhone(),
		UpdateAddressLine1:                  eventData.UpdateAddressLine1(),
		UpdateAddressLine2:                  eventData.UpdateAddressLine2(),
		UpdateAddressLine3:                  eventData.UpdateAddressLine3(),
		UpdateLocality:                      eventData.UpdateLocality(),
		UpdateCountry:                       eventData.UpdateCountry(),
		UpdateZip:                           eventData.UpdateZip(),
		UpdateLegalName:                     eventData.UpdateLegalName(),
		UpdateDomesticPaymentsBankInfo:      eventData.UpdateDomesticPaymentsBankInfo(),
		UpdateInternationalPaymentsBankInfo: eventData.UpdateInternationalPaymentsBankInfo(),
		UpdateVatNumber:                     eventData.UpdateVatNumber(),
		UpdateSendInvoicesFrom:              eventData.UpdateSendInvoicesFrom(),
		UpdateCanPayWithCard:                eventData.UpdateCanPayWithCard(),
		UpdateCanPayWithDirectDebitSEPA:     eventData.UpdateCanPayWithDirectDebitSEPA(),
		UpdateCanPayWithDirectDebitACH:      eventData.UpdateCanPayWithDirectDebitACH(),
		UpdateCanPayWithDirectDebitBacs:     eventData.UpdateCanPayWithDirectDebitBacs(),
		UpdateCanPayWithPigeon:              eventData.UpdateCanPayWithPigeon(),
	}
	err := h.repositories.Neo4jRepositories.TenantWriteRepository.UpdateTenantBillingProfile(ctx, tenantName, data)
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
		UpdatedAt:                  eventData.UpdatedAt,
		LogoUrl:                    eventData.LogoUrl,
		LogoRepositoryFileId:       eventData.LogoRepositoryFileId,
		InvoicingEnabled:           eventData.InvoicingEnabled,
		InvoicingPostpaid:          eventData.InvoicingPostpaid,
		BaseCurrency:               neo4jenum.DecodeCurrency(eventData.BaseCurrency),
		UpdateLogoUrl:              eventData.UpdateLogoUrl(),
		UpdateInvoicingEnabled:     eventData.UpdateInvoicingEnabled(),
		UpdateBaseCurrency:         eventData.UpdateBaseCurrency(),
		UpdateInvoicingPostpaid:    eventData.UpdateInvoicingPostpaid(),
		UpdateLogoRepositoryFileId: eventData.UpdateLogoRepositoryFileId(),
	}
	err := h.repositories.Neo4jRepositories.TenantWriteRepository.UpdateTenantSettings(ctx, tenantName, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return err
}
