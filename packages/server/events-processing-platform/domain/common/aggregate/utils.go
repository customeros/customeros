package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
)

func EnrichEventWithMetadata(event *eventstore.Event, span *opentracing.Span, tenant, userId string) {
	metadata := tracing.ExtractTextMapCarrier((*span).Context())
	metadata["tenant"] = tenant
	if userId != "" {
		metadata["user-id"] = userId
	}
	if err := event.SetMetadata(metadata); err != nil {
		tracing.TraceErr(*span, err)
	}
}
