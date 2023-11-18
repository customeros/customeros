package mapper

import (
	entityDashboard "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity/dashboard"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapDashboardNewCustomersData(newCustomersData *entityDashboard.DashboardNewCustomersData) *model.DashboardNewCustomers {
	if newCustomersData == nil {
		return nil
	}
	return &model.DashboardNewCustomers{
		ThisMonthCount:              newCustomersData.ThisMonthCount,
		ThisMonthIncreasePercentage: newCustomersData.ThisMonthIncreasePercentage,
		PerMonth:                    MapDashboardNewCustomersMonthData(newCustomersData.Months),
	}
}
func MapDashboardNewCustomersMonthData(months []*entityDashboard.DashboardNewCustomerMonthData) []*model.DashboardNewCustomersPerMonth {
	var result []*model.DashboardNewCustomersPerMonth
	for _, month := range months {
		result = append(result, &model.DashboardNewCustomersPerMonth{
			Month: month.Month,
			Count: month.Count,
		})
	}
	return result
}

func MapDashboardRetentionRateData(retentionRateData *entityDashboard.DashboardRetentionRateData) *model.DashboardRetentionRate {
	if retentionRateData == nil {
		return nil
	}
	return &model.DashboardRetentionRate{
		RetentionRate:      retentionRateData.RetentionRate,
		IncreasePercentage: retentionRateData.IncreasePercentage,
		PerMonth:           MapDashboardRetentionRatePerMonthData(retentionRateData.Months),
	}
}
func MapDashboardRetentionRatePerMonthData(months []*entityDashboard.DashboardRetentionRatePerMonthData) []*model.DashboardRetentionRatePerMonth {
	var result []*model.DashboardRetentionRatePerMonth
	for _, month := range months {
		result = append(result, &model.DashboardRetentionRatePerMonth{
			Month:      month.Month,
			RenewCount: month.RenewCount,
			ChurnCount: month.ChurnCount,
		})
	}
	return result
}

func MapDashboardRevenueAtRiskData(retentionRateData *entityDashboard.DashboardRevenueAtRiskData) *model.DashboardRevenueAtRisk {
	if retentionRateData == nil {
		return nil
	}
	return &model.DashboardRevenueAtRisk{
		HighConfidence: retentionRateData.HighConfidence,
		AtRisk:         retentionRateData.AtRisk,
	}
}
