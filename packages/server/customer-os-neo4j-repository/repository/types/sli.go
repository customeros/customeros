package types

import "time"

type SLI struct {
	Id          string
	BillingType string
	Price       float64
	Quantity    int64
	StartedAt   time.Time
	ContractId  string
}
