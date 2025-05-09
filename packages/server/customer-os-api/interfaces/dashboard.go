package cosapi_interfaces

import (
	"context"

	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

type DashboardService interface {
	GetDashboardViewOrganizationsData(ctx context.Context, requestDetails DashboardViewOrganizationsRequest) (*utils.Pagination, error)
	GetDashboardViewRenewalsData(ctx context.Context, requestDetails DashboardViewRenewalsRequest) (*utils.Pagination, error)
}

type DashboardViewOrganizationsRequest struct {
	Where *model.Filter
	Sort  *commonModel.SortBy
	Page  int
	Limit int
}

type DashboardViewRenewalsRequest struct {
	Where *model.Filter
	Sort  *commonModel.SortBy
	Page  int
	Limit int
}
