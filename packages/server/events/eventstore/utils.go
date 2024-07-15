package eventstore

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	baseEvent "github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"strings"
)

type EventMetadata struct {
	Tenant string `json:"tenant"`
	UserId string `json:"user-id"`
	App    string `json:"app"`
}

// TODO
func ToAggregateEvent(aggregate Aggregate, eventData baseEvent.BaseEvent) (Event, error) {
	if err := validator.GetValidator().Struct(eventData); err != nil {
		return Event{}, errors.Wrap(err, "failed to validate eventData")
	}

	event := NewBaseEvent(aggregate, eventData.EventName)
	if err := event.SetJsonData(&eventData); err != nil {
		return Event{}, errors.Wrap(err, "error setting json data for EmailCreateEvent")
	}
	return event, nil
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
	return (appSource == utils.AppSourceIntegrationApp || appSource == utils.AppSourceSyncCustomerOsData) && loggedInUserId == ""
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
	if strings.HasPrefix(aggregateID, string(aggregateType)+"-"+utils.StreamTempPrefix+"-"+tenant+"-") {
		return strings.ReplaceAll(aggregateID, string(aggregateType)+"-"+utils.StreamTempPrefix+"-"+tenant+"-", "")
	}
	return strings.ReplaceAll(aggregateID, string(aggregateType)+"-"+tenant+"-", "")
}

func GetTenantFromAggregate(aggregateID string, aggregateType AggregateType) string {
	if strings.HasPrefix(aggregateID, string(aggregateType)+"-"+utils.StreamTempPrefix+"-") {
		return strings.ReplaceAll(aggregateID, string(aggregateType)+"-"+utils.StreamTempPrefix+"-", "")
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
