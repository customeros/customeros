package invoicing_cycle

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	invoicingcyclepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoicing_cycle"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type CreateInvoicingCycleHandler interface {
	Handle(ctx context.Context, baseRequest events.BaseRequest, request *invoicingcyclepb.CreateInvoicingCycleTypeRequest) error
}

type createInvoicingCycleHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewCreateInvoicingCycleHandler(log logger.Logger, es eventstore.AggregateStore) CreateInvoicingCycleHandler {
	return &createInvoicingCycleHandler{log: log, es: es}
}

func (h *createInvoicingCycleHandler) Handle(ctx context.Context, baseRequest events.BaseRequest, request *invoicingcyclepb.CreateInvoicingCycleTypeRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateInvoicingCycleHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, baseRequest.Tenant, baseRequest.LoggedInUserId)
	tracing.LogObjectAsJson(span, "common", baseRequest)
	tracing.LogObjectAsJson(span, "request", request)

	invoicingCycleAggregate, err := LoadInvoicingCycleAggregate(ctx, h.es, baseRequest.Tenant, baseRequest.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	createEvent, err := NewInvoicingCycleCreateEvent(invoicingCycleAggregate, string(InvoicingCycleType(request.Type).StringValue()), utils.TimestampProtoToTimePtr(request.CreatedAt), baseRequest.SourceFields)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInvoicingCycleCreateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&createEvent, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	err = invoicingCycleAggregate.Apply(createEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Persist the changes to the event store
	if err = h.es.Save(ctx, invoicingCycleAggregate); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
