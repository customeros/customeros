package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type SLIProperty string

const (
	SLIPropertyPaused SLIProperty = "paused"
)

type ServiceLineItemEntity struct {
	ID               string
	Name             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	StartedAt        time.Time  // DateTime
	EndedAt          *time.Time // DateTime
	Canceled         bool
	Billed           enum.BilledType
	Price            float64
	Quantity         int64
	PreviousBilled   enum.BilledType
	PreviousPrice    float64
	PreviousQuantity int64
	Comments         string
	Source           DataSource
	SourceOfTruth    DataSource
	AppSource        string
	ParentID         string
	VatRate          float64
	Paused           bool

	DataLoaderKey
}

type ServiceLineItemEntities []ServiceLineItemEntity

func (sli ServiceLineItemEntity) IsEnded() bool {
	return sli.EndedAt != nil && sli.EndedAt.Before(utils.Now())
}

func (sli ServiceLineItemEntity) IsActiveAt(referenceTime time.Time) bool {
	return (sli.StartedAt.Equal(referenceTime) || sli.StartedAt.Before(referenceTime)) && (sli.EndedAt == nil || sli.EndedAt.After(referenceTime))
}

func (sli ServiceLineItemEntity) IsRecurrent() bool {
	return sli.Billed == enum.BilledTypeMonthly || sli.Billed == enum.BilledTypeAnnually || sli.Billed == enum.BilledTypeQuarterly
}

func (sli ServiceLineItemEntity) IsOneTime() bool {
	return sli.Billed == enum.BilledTypeOnce
}
