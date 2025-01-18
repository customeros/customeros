package cosapi_interfaces

import (
	"context"
	"time"

	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"

	entityDashboard "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity/dashboard"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

type DashboardService interface {
	GetDashboardViewOrganizationsData(ctx context.Context, requestDetails DashboardViewOrganizationsRequest) (*utils.Pagination, error)
	GetDashboardViewRenewalsData(ctx context.Context, requestDetails DashboardViewRenewalsRequest) (*utils.Pagination, error)

	GetDashboardCustomerMapData(ctx context.Context) ([]*entityDashboard.DashboardCustomerMapData, error)
	GetDashboardMRRPerCustomerData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardDashboardMRRPerCustomerData, error)
	GetDashboardGrossRevenueRetentionData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardGrossRevenueRetentionData, error)
	GetDashboardARRBreakdownData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardARRBreakdownData, error)
	GetDashboardRevenueAtRiskData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardRevenueAtRiskData, error)
	GetDashboardRetentionRateData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardRetentionRateData, error)
	GetDashboardNewCustomersData(ctx context.Context, start, end time.Time) (*entityDashboard.DashboardNewCustomersData, error)
	GetDashboardAverageTimeToOnboardPerMonth(ctx context.Context, start, end time.Time) (*model.DashboardTimeToOnboard, error)
	GetDashboardOnboardingCompletionPerMonth(ctx context.Context, start, end time.Time) (*model.DashboardOnboardingCompletion, error)
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
