package graph

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/customeros/customeros/packages/server/events-processing-platform/domain/invoice"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/constants"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/logger"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/tracing"
)

type InvoiceActionMetadata struct {
	Status        string  `json:"status"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	InvoiceNumber string  `json:"number"`
	InvoiceId     string  `json:"id"`
}

type InvoiceEventHandler struct {
	log         logger.Logger
	grpcClients *grpc_client.Clients
	neo4j       *neo4j_repository.Repositories
}

func NewInvoiceEventHandler(
	log logger.Logger,
	grpcClients *grpc_client.Clients,
	neo4j *neo4j_repository.Repositories,
) *InvoiceEventHandler {
	return &InvoiceEventHandler{
		log:         log,
		grpcClients: grpcClients,
		neo4j:       neo4j,
	}
}

func (h *InvoiceEventHandler) OnInvoiceVoidV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.OnInvoiceVoidV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceVoidEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	invoiceEntityBeforeVoid, err := h.getInvoice(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}

	err = h.neo4j.InvoiceWriteRepository.UpdateInvoice(ctx, nil, eventData.Tenant, invoiceId, neo4jrepository.InvoiceUpdateFields{
		UpdateStatus: true,
		Status:       neo4jenum.InvoiceStatusVoid,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while voiding invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}

	invoiceEntityAfterVoid, err := h.getInvoice(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while getting invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}
	h.createInvoiceAction(ctx, eventData.Tenant, invoiceEntityBeforeVoid.Status, *invoiceEntityAfterVoid)

	return nil
}

func (h *InvoiceEventHandler) OnInvoiceDeleteV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.OnInvoiceDeleteV1")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData invoice.InvoiceDeleteEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	invoiceId := invoice.GetInvoiceObjectID(evt.GetAggregateID(), eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, invoiceId)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	err := h.neo4j.InvoiceWriteRepository.DeleteInitializedInvoice(ctx, eventData.Tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while deleting invoice {%s}: {%s}", invoiceId, err.Error())
		return err
	}
	return err
}

func (h *InvoiceEventHandler) createInvoiceAction(ctx context.Context, tenant string, previousStatus neo4jenum.InvoiceStatus, invoiceEntity neo4jentity.InvoiceEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceEventHandler.createInvoiceAction")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("invoiceId", invoiceEntity.Id))
	span.LogFields(log.String("previousStatus", previousStatus.String()))
	span.LogFields(log.String("newStatus", invoiceEntity.Status.String()))
	span.LogFields(log.Bool("dryRun", invoiceEntity.DryRun))
	span.LogFields(log.Float64("totalAmount", invoiceEntity.TotalAmount))

	if previousStatus == invoiceEntity.Status {
		return
	}
	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) {
		span.LogFields(log.String("result", "dry run or total amount is 0"))
		return
	}

	metadata, err := utils.ToJson(InvoiceActionMetadata{
		Status:        invoiceEntity.Status.String(),
		Currency:      invoiceEntity.Currency.String(),
		Amount:        invoiceEntity.TotalAmount,
		InvoiceNumber: invoiceEntity.Number,
		InvoiceId:     invoiceEntity.Id,
	})

	actionType := enum.ActionNA
	message := ""
	switch invoiceEntity.Status {
	case neo4jenum.InvoiceStatusDue:
		message = "Invoice N° " + invoiceEntity.Number + " issued with an amount of " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)
		actionType = enum.ActionInvoiceIssued
	case neo4jenum.InvoiceStatusPaid:
		message = "Invoice N° " + invoiceEntity.Number + " paid in full: " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)
		actionType = enum.ActionInvoicePaid
	case neo4jenum.InvoiceStatusVoid:
		message = "Invoice N° " + invoiceEntity.Number + " voided"
		actionType = enum.ActionInvoiceVoided
	case neo4jenum.InvoiceStatusOverdue:
		message = "Invoice N° " + invoiceEntity.Number + " overdue"
		actionType = enum.ActionInvoiceOverdue
	}
	if actionType == enum.ActionNA {
		span.LogFields(log.String("result", "status not supported"))
		return
	}
	if invoiceEntity.Status == neo4jenum.InvoiceStatusDue {
		_, err = h.neo4j.ActionWriteRepository.MergeByActionType(ctx, nil, tenant, invoiceEntity.Id, model.INVOICE, actionType, message, metadata, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers)
	} else {
		_, err = h.neo4j.ActionWriteRepository.Create(ctx, tenant, invoiceEntity.Id, model.INVOICE, actionType, message, metadata, utils.Now(), constants.AppSourceEventProcessingPlatformSubscribers)
	}
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed creating invoice action for invoice %s: %s", invoiceEntity.Id, err.Error())
	}
}

func (h *InvoiceEventHandler) getInvoice(ctx context.Context, tenant, invoiceId string) (*neo4jentity.InvoiceEntity, error) {
	invoiceDbNode, err := h.neo4j.InvoiceReadRepository.GetInvoiceById(ctx, nil, tenant, invoiceId)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode), nil
}
