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

	SourceFields SourceFields

	DataloaderKey string
}

type OrderEntities []OrderEntity

func (OrderEntity) IsTimelineEvent() {
}

func (orderEntity *OrderEntity) SetDataloaderKey(key string) {
	orderEntity.DataloaderKey = key
}

func (orderEntity OrderEntity) GetDataloaderKey() string {
	return orderEntity.DataloaderKey
}

func (OrderEntity) TimelineEventLabel() string {
	return model.NodeLabelOrder
}

func (order OrderEntity) Labels(tenant string) []string {
	return []string{
		model.NodeLabelOrder,
		model.NodeLabelOrder + "_" + tenant,
	}
}
