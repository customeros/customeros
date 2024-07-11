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

type UpdateInvoicingCycleHandler interface {
	Handle(ctx context.Context, baseRequest events.BaseRequest, request *invoicingcyclepb.UpdateInvoicingCycleTypeRequest) error
}

type updateInvoicingCycleHandler struct {
	log logger.Logger
	es  eventstore.AggregateStore
}

func NewUpdateInvoicingCycleHandler(log logger.Logger, es eventstore.AggregateStore) UpdateInvoicingCycleHandler {
	return &updateInvoicingCycleHandler{log: log, es: es}
}

func (h *updateInvoicingCycleHandler) Handle(ctx context.Context, baseRequest events.BaseRequest, request *invoicingcyclepb.UpdateInvoicingCycleTypeRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UpdateInvoicingCycleHandler.Handle")
	defer span.Finish()
	tracing.SetCommandHandlerSpanTags(ctx, span, baseRequest.Tenant, baseRequest.LoggedInUserId)
	tracing.LogObjectAsJson(span, "common", baseRequest)
	tracing.LogObjectAsJson(span, "request", request)

	invoicingCycleAggregate, err := LoadInvoicingCycleAggregate(ctx, h.es, baseRequest.Tenant, baseRequest.ObjectID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	updateEvent, err := NewInvoicingCycleUpdateEvent(invoicingCycleAggregate, string(InvoicingCycleType(request.Type).StringValue()), utils.TimestampProtoToTimePtr(request.UpdatedAt), baseRequest.SourceFields)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInvoicingCycleCreateEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&updateEvent, span, eventstore.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	err = invoicingCycleAggregate.Apply(updateEvent)
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
