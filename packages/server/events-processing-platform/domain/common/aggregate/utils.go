package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
)

type EventMetadata struct {
	Tenant string `json:"tenant"`
	UserId string `json:"user-id"`
	App    string `json:"app"`
}

// Deprecated, use EnrichEventWithMetadataExtended instead
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

func EnrichEventWithMetadataExtended(event *eventstore.Event, span opentracing.Span, mtd EventMetadata) {
	metadata := tracing.ExtractTextMapCarrier(span.Context())
	metadata["tenant"] = mtd.Tenant
	if mtd.UserId != "" {
		metadata["user-id"] = mtd.UserId
	}
	if mtd.App != "" {
		metadata["app"] = mtd.App
	}
	if err := event.SetMetadata(metadata); err != nil {
		tracing.TraceErr(span, err)
	}
}

func AllowCheckIfEventIsRedundant(appSource, loggedInUserId string) bool {
	return (appSource == constants.AppSourceIntegrationApp || appSource == constants.AppSourceSyncCustomerOsData) && loggedInUserId == ""
}
