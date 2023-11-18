package entity

type DashboardNewCustomersData struct {
	ThisMonthCount              int
	ThisMonthIncreasePercentage float64
	Months                      []*DashboardNewCustomerMonthData
}
type DashboardNewCustomerMonthData struct {
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
