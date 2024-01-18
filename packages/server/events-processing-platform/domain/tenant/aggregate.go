package invoice

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/tenant/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	tenantpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/tenant"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

const (
	TenantAggregateType eventstore.AggregateType = "tenant"
)

type TenantAggregate struct {
	*aggregate.CommonIdAggregate
	TenantDetails *Tenant
}

func GetTenantName(aggregateID string) string {
	return strings.ReplaceAll(aggregateID, string(TenantAggregateType)+"-", "")
}

func LoadTenantAggregate(ctx context.Context, eventStore eventstore.AggregateStore, tenant string, options eventstore.LoadAggregateOptions) (*TenantAggregate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LoadTenantAggregate")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)

	tenantAggregate := NewTenantAggregate(tenant)

	err := aggregate.LoadAggregate(ctx, eventStore, tenantAggregate, options)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return tenantAggregate, nil
}

func NewTenantAggregate(tenant string) *TenantAggregate {
	tenantAggregate := TenantAggregate{}
	tenantAggregate.CommonIdAggregate = aggregate.NewCommonAggregateWithId(TenantAggregateType, tenant)
	tenantAggregate.SetWhen(tenantAggregate.When)
	tenantAggregate.TenantDetails = &Tenant{}
	tenantAggregate.Tenant = tenant

	return &tenantAggregate
}

func (a *TenantAggregate) HandleRequest(ctx context.Context, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *tenantpb.AddBillingProfileRequest:
		return a.AddBillingProfile(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *TenantAggregate) AddBillingProfile(ctx context.Context, request *tenantpb.AddBillingProfileRequest) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "TenantAggregate.AddBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())

	billingProfileId := uuid.New().String()

	addBillingProfileEvent, err := event.NewCreateTenantBillingProfileEvent(a, sourceFields, billingProfileId, request.Email, request.Phone, request.AddressLine1, request.AddressLine2, request.AddressLine3, request.LegalName, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "CreateTenantBillingProfileEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&addBillingProfileEvent, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return billingProfileId, a.Apply(addBillingProfileEvent)
}

func (a *TenantAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.TenantAddBillingProfileV1:
		return a.onAddBillingProfile(evt)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *TenantAggregate) onAddBillingProfile(evt eventstore.Event) error {
	var eventData event.CreateTenantBillingProfileEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.TenantDetails.HasBillingProfile(eventData.Id) {
		return nil
	}
	tenantBillingProfile := TenantBillingProfile{
		Id:           eventData.Id,
		CreatedAt:    eventData.CreatedAt,
		Email:        eventData.Email,
		Phone:        eventData.Phone,
		AddressLine1: eventData.AddressLine1,
		AddressLine2: eventData.AddressLine2,
		AddressLine3: eventData.AddressLine3,
		LegalName:    eventData.LegalName,
		SourceFields: eventData.SourceFields,
	}
	a.TenantDetails.BillingProfiles = append(a.TenantDetails.BillingProfiles, &tenantBillingProfile)

	return nil
}
