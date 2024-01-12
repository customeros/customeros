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

type InvoiceFillHandler interface {
	Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *invoicepb.FillInvoiceRequest) error
}

type invoiceFillHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewInvoiceFillHandler(log logger.Logger, es eventstore.AggregateStore) InvoiceFillHandler {
	return &invoiceFillHandler{log: log, es: es}
}

func (h *invoiceFillHandler) Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *invoicepb.FillInvoiceRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceFillHandler.Handle")
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

	fillEvent, err := NewInvoiceFillEvent(invoiceAggregate, utils.TimestampProtoToTimePtr(request.UpdatedAt), baseRequest.SourceFields, request)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceFillEvent")
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

	// Persist the changes to the event store
	if err = h.es.Save(ctx, invoiceAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
