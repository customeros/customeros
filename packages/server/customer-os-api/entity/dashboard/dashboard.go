package entity

import (
	"time"
)

type DashboardNewCustomersData struct {
	ThisMonthCount              int
	ThisMonthIncreasePercentage float64
	Months                      []*DashboardNewCustomerMonthData
}
type DashboardNewCustomerMonthData struct {
	Year  int
	Month int
	Count int
}

type DashboardRetentionRateData struct {
	RetentionRate      int
	IncreasePercentage float64
	Months             []*DashboardRetentionRatePerMonthData
}
type DashboardRetentionRatePerMonthData struct {
	Month         int
	RenewCount    int
	ChurnCount    int
	RetentionRate int
}

type DashboardRevenueAtRiskData struct {
	HighConfidence float64
	AtRisk         float64
}

type DashboardARRBreakdownData struct {
	ArrBreakdown       float64
	IncreasePercentage float64
	Months             []*DashboardARRBreakdownPerMonthData
}
type DashboardARRBreakdownPerMonthData struct {
	Month           int
	NewlyContracted int
	Renewals        int
	Upsells         int
	Downgrades      int
	Cancellations   int
	Churned         int
}

type DashboardGrossRevenueRetentionData struct {
	GrossRevenueRetention float64
	IncreasePercentage    float64
	Months                []*DashboardGrossRevenueRetentionPerMonthData
}
type DashboardGrossRevenueRetentionPerMonthData struct {
	Month      int
	Percentage float64
}

type DashboardDashboardMRRPerCustomerData struct {
	MrrPerCustomer     float64 `json:"mrrPerCustomer"`
	IncreasePercentage float64
	Months             []*DashboardDashboardMRRPerCustomerPerMonthData
}
type DashboardDashboardMRRPerCustomerPerMonthData struct {
	Year  int
	Month int
	Value float64
}

type DashboardCustomerMapState string

const (
	DashboardCustomerMapStateOk      DashboardCustomerMapState = "OK"
	DashboardCustomerMapStateAtRisk  DashboardCustomerMapState = "AT_RISK"
	DashboardCustomerMapStateChurned DashboardCustomerMapState = "CHURNED"
)

var DashboardCustomerMapStates = []DashboardCustomerMapState{
	DashboardCustomerMapStateOk,
	DashboardCustomerMapStateAtRisk,
	DashboardCustomerMapStateChurned,
}

type DashboardCustomerMapData struct {
	OrganizationId     string
	State              DashboardCustomerMapState
	Arr                float64
	ContractSignedDate time.Time
}
