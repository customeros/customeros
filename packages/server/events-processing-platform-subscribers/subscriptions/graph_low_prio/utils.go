package graph_low_prio

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/opentracing/opentracing-go"
)

func setEventSpanTagsAndLogFields(span opentracing.Span, evt eventstore.Event) {
	span.SetTag(tracing.SpanTagComponent, constants.ComponentSubscriptionGraph)
	span.SetTag(tracing.SpanTagAggregateId, evt.GetAggregateID())
	tracing.LogObjectAsJson(span, "event", evt)
}
