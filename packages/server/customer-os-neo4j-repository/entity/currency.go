package entity

import (
	"time"
)

type CurrencyEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Name   string
	Symbol string

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}
