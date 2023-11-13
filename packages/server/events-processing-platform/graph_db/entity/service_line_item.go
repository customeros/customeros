package entity

import (
	"time"
)

type ServiceLineItemEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	Description   string
	Billed        string
	Price         float64
	Licenses      int64
}
