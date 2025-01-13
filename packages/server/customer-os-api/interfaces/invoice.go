package cosapi_interfaces

import (
	"context"

	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

type InvoiceService interface {
	CountInvoices(ctx context.Context, tenant, organizationId string, where *model.Filter) (int64, error)
	GetInvoices(ctx context.Context, organizationId string, page, limit int, where *model.Filter, sortBy []*commonModel.SortBy) (*utils.Pagination, error)
	UpdateInvoice(ctx context.Context, input model.InvoiceUpdateInput) error
}
