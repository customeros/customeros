package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type OrderEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time

	ConfirmedAt *time.Time
	PaidAt      *time.Time
	FulfilledAt *time.Time
	CancelledAt *time.Time

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	DataloaderKey string
}

type OrderEntities []OrderEntity

func (order OrderEntity) Labels(tenant string) []string {
	return []string{
		model.NodeLabelOrder,
		model.NodeLabelOrder + "_" + tenant,
	}
}
