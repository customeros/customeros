package eventstore

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"strings"
)

type EventMetadata struct {
	Tenant string `json:"tenant"`
	UserId string `json:"user-id"`
	App    string `json:"app"`
}

// Deprecated, use EnrichEventWithMetadataExtended instead
func EnrichEventWithMetadata(event *Event, span *opentracing.Span, tenant, userId string) {
	metadata := tracing.ExtractTextMapCarrier((*span).Context())
	metadata["tenant"] = tenant
	if userId != "" {
		metadata["user-id"] = userId
	}
	if err := event.SetMetadata(metadata); err != nil {
		tracing.TraceErr(*span, err)
	}
}

func EnrichEventWithMetadataExtended(event *Event, span opentracing.Span, mtd EventMetadata) {
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

func AllowCheckForNoChanges(appSource, loggedInUserId string) bool {
	return (appSource == events.AppSourceIntegrationApp || appSource == events.AppSourceSyncCustomerOsData) && loggedInUserId == ""
}

func LoadAggregate(ctx context.Context, eventStore AggregateStore, agg Aggregate, options LoadAggregateOptions) error {
	err := eventStore.Exists(ctx, agg.GetID())
	if err != nil {
		if !errors.Is(err, ErrAggregateNotFound) {
			return err
		} else {
			return nil
		}
	}

	if options.SkipLoadEvents {
		if err = eventStore.LoadVersion(ctx, agg); err != nil {
			return err
		}
	} else {
		if err = eventStore.Load(ctx, agg); err != nil {
			return err
		}
	}

	return nil
}

func GetAggregateObjectID(aggregateID, tenant string, aggregateType AggregateType) string {
	if tenant == "" {
		return getAggregateObjectUUID(aggregateID)
	}
	if strings.HasPrefix(aggregateID, string(aggregateType)+"-"+events.StreamTempPrefix+"-"+tenant+"-") {
		return strings.ReplaceAll(aggregateID, string(aggregateType)+"-"+events.StreamTempPrefix+"-"+tenant+"-", "")
	}
	return strings.ReplaceAll(aggregateID, string(aggregateType)+"-"+tenant+"-", "")
}

func GetTenantFromAggregate(aggregateID string, aggregateType AggregateType) string {
	if strings.HasPrefix(aggregateID, string(aggregateType)+"-"+events.StreamTempPrefix+"-") {
		return strings.ReplaceAll(aggregateID, string(aggregateType)+"-"+events.StreamTempPrefix+"-", "")
	}

	var1 := strings.ReplaceAll(aggregateID, string(aggregateType)+"-", "")
	var2 := strings.Split(var1, "-")
	return var2[0]
}

// use this method when tenant is not known
func getAggregateObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}
