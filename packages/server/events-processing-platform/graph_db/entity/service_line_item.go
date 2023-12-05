package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type ServiceLineItemEntity struct {
	Id            string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StartedAt     time.Time
	EndedAt       *time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	Name          string
	Billed        string
	Price         float64
	Quantity      int64
	Comments      string
	ParentId      string
	IsCanceled    bool
}

type ServiceLineItemEntities []ServiceLineItemEntity

func (sli ServiceLineItemEntity) IsEnded() bool {
	return sli.EndedAt != nil && sli.EndedAt.Before(utils.Now())
}
