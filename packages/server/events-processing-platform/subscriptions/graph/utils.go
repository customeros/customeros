package graph

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"math/rand"
	"time"
)

var (
	charset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
)

func generateNewRandomCustomerOsId() string {
	customerOsID := "C-" + generateRandomStringFromCharset(3) + "-" + generateRandomStringFromCharset(3)
	return customerOsID
}

func generateRandomStringFromCharset(length int) string {
	// Create a new source based on the current time's Unix timestamp (in nanoseconds)
	source := rand.NewSource(time.Now().UnixNano())
	// Initialize a random number generator (RNG) with the source
	rng := rand.New(source)

	var output string
	for i := 0; i < length; i++ {
		randChar := charset[rng.Intn(len(charset))]
		output += string(randChar)
	}
	return output
}

func setCommonSpanTagsAndLogFields(span opentracing.Span, evt eventstore.Event) {
	span.SetTag(tracing.SpanTagComponent, constants.ComponentSubscriptionGraph)
	span.SetTag(tracing.SpanTagAggregateId, evt.GetAggregateID())
	span.LogFields(log.Object("event", evt))
}
