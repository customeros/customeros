package cosapi_interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

type BillableService interface {
	GetBillableDetails(ctx context.Context) (*model.TenantBillableInfo, error)
}
