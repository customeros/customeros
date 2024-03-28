package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type OfferingEntity struct {
	Id                    string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Source                DataSource
	SourceOfTruth         DataSource
	AppSource             string
	Name                  string
	Active                bool
	Type                  enum.OfferingType
	PricingModel          enum.PricingModel
	PricingPeriodInMonths int64
	Currency              enum.Currency
	Price                 float64
	PriceCalculated       bool
	Conditional           bool
	Taxable               bool
	PriceCalculation      PriceCalculation
	Conditionals          Conditionals
}

type PriceCalculation struct {
	Type                   enum.PriceCalculationType
	RevenueSharePercentage float64
}

type Conditionals struct {
	MinimumChargePeriod enum.ChargePeriod
	MinimumChargeAmount float64
}

type OfferingEntities []OfferingEntity
