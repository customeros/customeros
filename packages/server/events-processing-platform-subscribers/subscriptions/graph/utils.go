package graph

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/constants"
	"github.com/customeros/customeros/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
)

var (
	charset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
)

func generateNewRandomCustomerOsId() string {
	customerOsID := "C-" + utils.GenerateRandomStringFromCharset(3, charset) + "-" + utils.GenerateRandomStringFromCharset(3, charset)
	return customerOsID
}

func setEventSpanTagsAndLogFields(span opentracing.Span, evt eventstore.Event) {
	span.SetTag(tracing.SpanTagComponent, constants.ComponentSubscriptionGraph)
	span.SetTag(tracing.SpanTagAggregateId, evt.GetAggregateID())
	tracing.LogObjectAsJson(span, "event", evt)
}
