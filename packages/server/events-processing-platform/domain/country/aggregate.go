package country

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	countrypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/country"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const (
	CountryAggregateType eventstore.AggregateType = "country"
)

type countryAggregate struct {
	*eventstore.CommonIdAggregate
	Country *Country
}

func GetCountryObjectID(aggregateID string) string {
	return getCountryObjectUUID(aggregateID)
}

func getCountryObjectUUID(aggregateID string) string {
	parts := strings.Split(aggregateID, "-")
	fullUUID := parts[len(parts)-5] + "-" + parts[len(parts)-4] + "-" + parts[len(parts)-3] + "-" + parts[len(parts)-2] + "-" + parts[len(parts)-1]
	return fullUUID
}

func NewCountryAggregateWithID(id string) *countryAggregate {
	countryAggregate := countryAggregate{}
	countryAggregate.CommonIdAggregate = eventstore.NewCommonAggregateWithId(CountryAggregateType, id)
	countryAggregate.SetWhen(countryAggregate.When)
	countryAggregate.Country = &Country{}

	return &countryAggregate
}

func (a *countryAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CountryAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *countrypb.CreateCountryRequest:
		return nil, a.CreateCountryRequest(ctx, r)
	default:
		return nil, nil
	}
}

func (a *countryAggregate) CreateCountryRequest(ctx context.Context, request *countrypb.CreateCountryRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "CountryAggregate.CreateCountryRequest")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	createEvent, err := NewCountryCreateEvent(a, request.Name, request.CodeA2, request.CodeA3, request.PhoneCode, time.Now().UTC())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "CountryAggregate")
	}

	eventstore.EnrichEventWithMetadataExtended(&createEvent, span, eventstore.EventMetadata{
		UserId: request.LoggedInUserId,
		App:    request.SourceFields.AppSource,
	})

	streamMetadata := esdb.StreamMetadata{}
	streamMetadata.SetMaxAge(time.Duration(utils.StreamMetadataMaxAgeSecondsExtended) * time.Second)
	a.SetStreamMetadata(&streamMetadata)

	return a.Apply(createEvent)
}

func (a *countryAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case CountryCreateV1:
		return a.onCountryCreate(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), utils.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *countryAggregate) onCountryCreate(evt eventstore.Event) error {
	var eventData CountryCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Country.CreatedAt = eventData.CreatedAt
	a.Country.Name = eventData.Name
	a.Country.CodeA2 = eventData.CodeA2
	a.Country.CodeA3 = eventData.CodeA3
	a.Country.PhoneCode = eventData.PhoneCode

	return nil
}
