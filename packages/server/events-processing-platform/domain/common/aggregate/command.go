package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type TenantEventInterface interface {
	GetTenant() string
}

func CreateEvent(ctx context.Context, operationName string,
	aggregate eventstore.Aggregate, command TenantEventInterface,
	createEvent func() (eventstore.Event, error)) error {
	span, _ := opentracing.StartSpanFromContext(ctx, operationName)
	defer span.Finish()
	span.LogFields(log.String("Tenant", command.GetTenant()), log.String("Id", aggregate.GetID()))

	event, err := createEvent()
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, operationName)
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, operationName)
	}

	return aggregate.Apply(event)
}
