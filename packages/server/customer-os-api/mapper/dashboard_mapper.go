package mapper

import (
	"fmt"
	entityDashboard "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity/dashboard"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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
			Year:  month.Year,
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
	retentionRateModel := model.DashboardRetentionRate{
		RetentionRate:           utils.TruncateFloat64(retentionRateData.RetentionRate, 1),
		IncreasePercentageValue: utils.TruncateFloat64(retentionRateData.IncreasePercentage, 1),
		IncreasePercentage:      fmt.Sprintf("%.1f", utils.TruncateFloat64(retentionRateData.IncreasePercentage, 1)),
		PerMonth:                MapDashboardRetentionRatePerMonthData(retentionRateData.Months),
	}
	return &retentionRateModel
}

func MapDashboardRetentionRatePerMonthData(months []*entityDashboard.DashboardRetentionRatePerMonthData) []*model.DashboardRetentionRatePerMonth {
	var result []*model.DashboardRetentionRatePerMonth
	for _, month := range months {
		result = append(result, &model.DashboardRetentionRatePerMonth{
			Year:       month.Year,
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
		HighConfidence: utils.TruncateFloat64(retentionRateData.HighConfidence, 1),
		AtRisk:         utils.TruncateFloat64(retentionRateData.AtRisk, 1),
	}
}

func MapDashboardARRBreakdownData(retentionRateData *entityDashboard.DashboardARRBreakdownData) *model.DashboardARRBreakdown {
	if retentionRateData == nil {
		return nil
	}
	return &model.DashboardARRBreakdown{
		ArrBreakdown:       utils.TruncateFloat64(retentionRateData.ArrBreakdown, 1),
		IncreasePercentage: retentionRateData.IncreasePercentage,
		PerMonth:           MapDashboardARRBreakdownPerMonthData(retentionRateData.Months),
	}
}

func MapDashboardARRBreakdownPerMonthData(months []*entityDashboard.DashboardARRBreakdownPerMonthData) []*model.DashboardARRBreakdownPerMonth {
	var result []*model.DashboardARRBreakdownPerMonth
	for _, month := range months {
		result = append(result, &model.DashboardARRBreakdownPerMonth{
			Year:            month.Year,
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
		GrossRevenueRetention:   utils.TruncateFloat64(grossRevenueRetentionData.GrossRevenueRetention, 1),
		IncreasePercentageValue: utils.TruncateFloat64(grossRevenueRetentionData.IncreasePercentage, 1),
		IncreasePercentage:      fmt.Sprintf("%.1f", utils.TruncateFloat64(grossRevenueRetentionData.IncreasePercentage, 1)),
		PerMonth:                MapDashboardGrossRevenueRetentionPerMonthData(grossRevenueRetentionData.Months),
	}
}
func MapDashboardGrossRevenueRetentionPerMonthData(months []*entityDashboard.DashboardGrossRevenueRetentionPerMonthData) []*model.DashboardGrossRevenueRetentionPerMonth {
	var result []*model.DashboardGrossRevenueRetentionPerMonth
	for _, month := range months {
		result = append(result, &model.DashboardGrossRevenueRetentionPerMonth{
			Year:       month.Year,
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
			Year:  month.Year,
			Month: month.Month,
			Value: month.Value,
		})
	}
	return result
}

func MapDashboardCustomerMapDataList(grossRevenueRetentionData []*entityDashboard.DashboardCustomerMapData) []*model.DashboardCustomerMap {
	var result []*model.DashboardCustomerMap
	for _, month := range grossRevenueRetentionData {
		result = append(result, MapDashboardCustomerMapData(month))
	}
	return result
}
func MapDashboardCustomerMapData(dashboardCustomerMapData *entityDashboard.DashboardCustomerMapData) *model.DashboardCustomerMap {
	if dashboardCustomerMapData == nil {
		return nil
	}
	return &model.DashboardCustomerMap{
		OrganizationID:     dashboardCustomerMapData.OrganizationId,
		State:              MapDashboardCustomerMapState(dashboardCustomerMapData.State),
		Arr:                dashboardCustomerMapData.Arr,
		ContractSignedDate: dashboardCustomerMapData.ContractSignedDate,
	}
}
func MapDashboardCustomerMapState(state entityDashboard.DashboardCustomerMapState) model.DashboardCustomerMapState {
	switch state {
	case entityDashboard.DashboardCustomerMapStateOk:
		return model.DashboardCustomerMapStateOk
	case entityDashboard.DashboardCustomerMapStateAtRisk:
		return model.DashboardCustomerMapStateAtRisk
	case entityDashboard.DashboardCustomerMapStateChurned:
		return model.DashboardCustomerMapStateChurned
	case entityDashboard.DashboardCustomerMapStateHighRisk:
		return model.DashboardCustomerMapStateHighRisk
	case entityDashboard.DashboardCustomerMapStateMediumRisk:
		return model.DashboardCustomerMapStateMediumRisk
	default:
		return ""
	}
}
