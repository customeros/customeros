package invoice

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type InvoicePdfGeneratedHandler interface {
	Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *invoicepb.PdfGeneratedInvoiceRequest) error
}

type invoicePdfGeneratedHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewInvoicePdfGeneratedHandler(log logger.Logger, es eventstore.AggregateStore) InvoicePdfGeneratedHandler {
	return &invoicePdfGeneratedHandler{log: log, es: es}
}

func (h *invoicePdfGeneratedHandler) Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *invoicepb.PdfGeneratedInvoiceRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoicePdfGeneratedHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, baseRequest.Tenant, baseRequest.LoggedInUserId)
	tracing.LogObjectAsJson(span, "common", baseRequest)
	tracing.LogObjectAsJson(span, "request", request)

	invoiceAggregate, err := LoadInvoiceAggregate(ctx, h.es, baseRequest.Tenant, baseRequest.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if eventstore.IsAggregateNotFound(invoiceAggregate) {
		tracing.TraceErr(span, eventstore.ErrAggregateNotFound)
		return eventstore.ErrAggregateNotFound
	}

	fillEvent, err := NewInvoicePdfGeneratedEvent(invoiceAggregate, utils.TimestampProtoToTimePtr(request.UpdatedAt), baseRequest.SourceFields, request)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoicePdfGeneratedEvent")
	}
	commonAggregate.EnrichEventWithMetadataExtended(&fillEvent, span, commonAggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	err = invoiceAggregate.Apply(fillEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err = h.es.Save(ctx, invoiceAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
