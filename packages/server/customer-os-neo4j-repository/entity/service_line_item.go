package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type ServiceLineItemEntity struct {
	DataLoaderKey
	ID               string
	Name             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	StartedAt        time.Time
	EndedAt          *time.Time
	IsCanceled       bool
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
}

type ServiceLineItemEntities []ServiceLineItemEntity

func (sli ServiceLineItemEntity) IsEnded() bool {
	return sli.EndedAt != nil && sli.EndedAt.Before(utils.Now())
}

func (sli ServiceLineItemEntity) IsParent() bool {
	return sli.ParentID == sli.ID
}

func (sli ServiceLineItemEntity) IsActiveAt(referenceTime time.Time) bool {
	return (sli.StartedAt.Equal(referenceTime) || sli.StartedAt.Before(referenceTime)) && (sli.EndedAt == nil || sli.EndedAt.After(referenceTime))
}
