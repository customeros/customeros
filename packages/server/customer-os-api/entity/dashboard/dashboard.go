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
