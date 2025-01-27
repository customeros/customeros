package aggregate

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	organizationpb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/customeros/customeros/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type OrganizationTempAggregate struct {
	*eventstore.CommonTenantIdTempAggregate
}

func NewOrganizationTempAggregateWithTenantAndID(tenant, id string) *OrganizationTempAggregate {
	organizationTempAggregate := OrganizationTempAggregate{}
	organizationTempAggregate.CommonTenantIdTempAggregate = eventstore.NewCommonTempAggregateWithTenantAndId(OrganizationAggregateType, tenant, id)
	organizationTempAggregate.Tenant = tenant

	return &organizationTempAggregate
}

func (a *OrganizationTempAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationTempAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch request.(type) {
	case *organizationpb.OrganizationIdGrpcRequest:
		rpc := params["rpc"]
		if rpc == nil {
			tracing.TraceErr(span, errors.New("rpc is nil"))
			return nil, errors.New("rpc is nil")
		}
		return nil, errors.New("invalid rpc")
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}
