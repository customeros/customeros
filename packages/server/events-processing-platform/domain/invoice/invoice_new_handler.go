package invoice

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonAggregate "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type InvoiceNewHandler interface {
	Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *invoicepb.NewInvoiceRequest) error
}

type invoiceNewHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewInvoiceNewHandler(log logger.Logger, es eventstore.AggregateStore) InvoiceNewHandler {
	return &invoiceNewHandler{log: log, es: es}
}

func (h *invoiceNewHandler) Handle(ctx context.Context, baseRequest eventstore.BaseRequest, request *invoicepb.NewInvoiceRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceNewHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, baseRequest.Tenant, baseRequest.LoggedInUserId)
	tracing.LogObjectAsJson(span, "common", baseRequest)
	tracing.LogObjectAsJson(span, "request", request)

	invoiceAggregate, err := LoadInvoiceAggregate(ctx, h.es, baseRequest.Tenant, baseRequest.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	date := utils.TimestampProtoToTime(request.Date).UTC()
	dueDate := utils.TimestampProtoToTime(request.Date).UTC()
	createdAt := utils.TimestampProtoToTime(request.CreatedAt).UTC()

	//TODO generate invoice number unique in tenant
	number := uuid.New().String()

	createEvent, err := NewInvoiceNewEvent(invoiceAggregate, request.ContractId, request.DryRun, number, date, dueDate, createdAt, baseRequest.SourceFields)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceNewEvent")
	}
	commonAggregate.EnrichEventWithMetadataExtended(&createEvent, span, commonAggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	err = invoiceAggregate.Apply(createEvent)
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
