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
