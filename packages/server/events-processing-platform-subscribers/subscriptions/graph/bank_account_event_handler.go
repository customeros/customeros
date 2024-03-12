package graph

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/helper"
	tenant "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type BankAccountEventHandler struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewBankAccountEventHandler(log logger.Logger, repositories *repository.Repositories) *BankAccountEventHandler {
	return &BankAccountEventHandler{
		log:          log,
		repositories: repositories,
	}
}

func (h *BankAccountEventHandler) OnAddBankAccountV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountEventHandler.OnAddBankAccountV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.TenantBankAccountCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	tenantName := tenant.GetTenantName(evt.GetAggregateID())
	span.SetTag(tracing.SpanTagEntityId, tenantName)
	span.SetTag(tracing.SpanTagEntityId, eventData.Id)

	data := neo4jrepository.BankAccountCreateFields{
		Id:        eventData.Id,
		CreatedAt: eventData.CreatedAt,
		SourceFields: neo4jmodel.Source{
			Source:    helper.GetSource(eventData.SourceFields.Source),
			AppSource: helper.GetAppSource(eventData.SourceFields.AppSource),
		},
		BankName:            eventData.BankName,
		BankTransferEnabled: eventData.BankTransferEnabled,
		AllowInternational:  eventData.AllowInternational,
		Currency:            neo4jenum.DecodeCurrency(eventData.Currency),
		AccountNumber:       eventData.AccountNumber,
		SortCode:            eventData.SortCode,
		Iban:                eventData.Iban,
		Bic:                 eventData.Bic,
		RoutingNumber:       eventData.RoutingNumber,
	}
	err := h.repositories.Neo4jRepositories.BankAccountWriteRepository.CreateBankAccount(ctx, tenantName, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return err
}

func (h *BankAccountEventHandler) OnUpdateBankAccountV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountEventHandler.OnUpdateBankAccountV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.TenantBankAccountUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	tenantName := tenant.GetTenantName(evt.GetAggregateID())
	span.SetTag(tracing.SpanTagEntityId, tenantName)
	span.SetTag(tracing.SpanTagEntityId, eventData.Id)

	data := neo4jrepository.BankAccountUpdateFields{
		Id:                        eventData.Id,
		UpdatedAt:                 eventData.UpdatedAt,
		BankName:                  eventData.BankName,
		BankTransferEnabled:       eventData.BankTransferEnabled,
		AllowInternational:        eventData.AllowInternational,
		Currency:                  neo4jenum.DecodeCurrency(eventData.Currency),
		AccountNumber:             eventData.AccountNumber,
		SortCode:                  eventData.SortCode,
		Iban:                      eventData.Iban,
		Bic:                       eventData.Bic,
		RoutingNumber:             eventData.RoutingNumber,
		UpdateBankName:            eventData.UpdateBankName(),
		UpdateBankTransferEnabled: eventData.UpdateBankTransferEnabled(),
		UpdateAllowInternational:  eventData.UpdateAllowInternational(),
		UpdateCurrency:            eventData.UpdateCurrency(),
		UpdateAccountNumber:       eventData.UpdateAccountNumber(),
		UpdateSortCode:            eventData.UpdateSortCode(),
		UpdateIban:                eventData.UpdateIban(),
		UpdateBic:                 eventData.UpdateBic(),
		UpdateRoutingNumber:       eventData.UpdateRoutingNumber(),
	}
	err := h.repositories.Neo4jRepositories.BankAccountWriteRepository.UpdateBankAccount(ctx, tenantName, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return err
}

func (h *BankAccountEventHandler) OnDeleteBankAccountV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountEventHandler.OnDeleteBankAccountV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData event.TenantBankAccountDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	tenantName := tenant.GetTenantName(evt.GetAggregateID())
	span.SetTag(tracing.SpanTagEntityId, tenantName)

	err := h.repositories.Neo4jRepositories.BankAccountWriteRepository.DeleteBankAccount(ctx, tenantName, eventData.Id)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
