package offering

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	offeringpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/offering"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"strings"
)

const OfferingAggregateType = "offering"

type OfferingAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Offering *Offering
}

func GetOfferingObjectID(aggregateID string, tenant string) string {
	return aggregate.GetAggregateObjectID(aggregateID, tenant, OfferingAggregateType)
}

func LoadOfferingAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant, objectID string, opts eventstore.LoadAggregateOptions) (*OfferingAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadOfferingAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("ObjectID", objectID))

	OfferingAggregate := NewOfferingAggregateWithTenantAndID(tenant, objectID)

	err := aggregate.LoadAggregate(ctx, eventStore, OfferingAggregate, opts)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return OfferingAggregate, nil
}

func (a *OfferingAggregate) HandleRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OfferingAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *offeringpb.CreateOfferingGrpcRequest:
		return nil, a.CreateOffering(ctx, r)
	case *offeringpb.UpdateOfferingGrpcRequest:
		return nil, nil //TODO
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func NewOfferingAggregateWithTenantAndID(tenant, id string) *OfferingAggregate {
	OfferingAggregate := OfferingAggregate{}
	OfferingAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(OfferingAggregateType, tenant, id)
	OfferingAggregate.SetWhen(OfferingAggregate.When)
	OfferingAggregate.Offering = &Offering{}
	OfferingAggregate.Tenant = tenant

	return &OfferingAggregate
}

func (a *OfferingAggregate) When(event eventstore.Event) error {
	switch event.GetEventType() {
	case OfferingCreateV1:
		return a.whenOfferingCreate(event)
	default:
		if strings.HasPrefix(event.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *OfferingAggregate) whenOfferingCreate(evt eventstore.Event) error {

	return nil
}

func (a *OfferingAggregate) CreateOffering(ctx context.Context, request *offeringpb.CreateOfferingGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OfferingAggregate.CreateOffering")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	//createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	//sourceFields := commonmodel.Source{}
	//sourceFields.FromGrpc(request.SourceFields)
	//
	//createEvent, err := NewOfferingCreateEvent(
	//	a,
	//	request.Content,
	//	request.LoggedInUserId,
	//	request.OrganizationId,
	//	request.Dismissed,
	//	createdAtNotNil,
	//	dueDateNotNil,
	//	sourceFields,
	//)
	//if err != nil {
	//	tracing.TraceErr(span, err)
	//	return errors.Wrap(err, "NewOfferingCreateEvent")
	//}
	//aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
	//	Tenant: request.Tenant,
	//	UserId: request.LoggedInUserId,
	//	App:    request.SourceFields.AppSource,
	//})

	//return a.Apply(createEvent)
	return nil
}
