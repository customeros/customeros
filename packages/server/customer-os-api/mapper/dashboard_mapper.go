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

func MapDashboardARRBreakdownData(retentionRateData *entityDashboard.DashboardARRBreakdownData) *model.DashboardARRBreakdown {
	if retentionRateData == nil {
		return nil
	}
	return &model.DashboardARRBreakdown{
		ArrBreakdown:       retentionRateData.ArrBreakdown,
		IncreasePercentage: retentionRateData.IncreasePercentage,
		PerMonth:           MapDashboardARRBreakdownPerMonthData(retentionRateData.Months),
	}
}
func MapDashboardARRBreakdownPerMonthData(months []*entityDashboard.DashboardARRBreakdownPerMonthData) []*model.DashboardARRBreakdownPerMonth {
	var result []*model.DashboardARRBreakdownPerMonth
	for _, month := range months {
		result = append(result, &model.DashboardARRBreakdownPerMonth{
			Month:           month.Month,
			NewlyContracted: month.NewlyContracted,
			Renewals:        month.Renewals,
			Upsells:         month.Upsells,
			Downgrades:      month.Downgrades,
			Cancellations:   month.Cancellations,
			Churned:         month.Churned,
		})
	}
	return result
}

func MapDashboardGrossRevenueRetentionData(grossRevenueRetentionData *entityDashboard.DashboardGrossRevenueRetentionData) *model.DashboardGrossRevenueRetention {
	if grossRevenueRetentionData == nil {
		return nil
	}
	return &model.DashboardGrossRevenueRetention{
		GrossRevenueRetention: grossRevenueRetentionData.GrossRevenueRetention,
		IncreasePercentage:    grossRevenueRetentionData.IncreasePercentage,
		PerMonth:              MapDashboardGrossRevenueRetentionPerMonthData(grossRevenueRetentionData.Months),
	}
}
func MapDashboardGrossRevenueRetentionPerMonthData(months []*entityDashboard.DashboardGrossRevenueRetentionPerMonthData) []*model.DashboardGrossRevenueRetentionPerMonth {
	var result []*model.DashboardGrossRevenueRetentionPerMonth
	for _, month := range months {
		result = append(result, &model.DashboardGrossRevenueRetentionPerMonth{
			Month:      month.Month,
			Percentage: month.Percentage,
		})
	}
	return result
}

func MapDashboardMRRPerCustomerData(grossRevenueRetentionData *entityDashboard.DashboardDashboardMRRPerCustomerData) *model.DashboardMRRPerCustomer {
	if grossRevenueRetentionData == nil {
		return nil
	}
	return &model.DashboardMRRPerCustomer{
		MrrPerCustomer:     grossRevenueRetentionData.MrrPerCustomer,
		IncreasePercentage: grossRevenueRetentionData.IncreasePercentage,
		PerMonth:           MapDashboardMRRPerCustomerPerMonthData(grossRevenueRetentionData.Months),
	}
}
func MapDashboardMRRPerCustomerPerMonthData(months []*entityDashboard.DashboardDashboardMRRPerCustomerPerMonthData) []*model.DashboardMRRPerCustomerPerMonth {
	var result []*model.DashboardMRRPerCustomerPerMonth
	for _, month := range months {
		result = append(result, &model.DashboardMRRPerCustomerPerMonth{
			Month: month.Month,
			Value: month.Value,
		})
	}
	return result
}
