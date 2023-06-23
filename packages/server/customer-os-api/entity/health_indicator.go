package entity

import (
	"time"
)

type HealthIndicatorEntity struct {
	Id        string
	Name      string
	Order     int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Source    DataSource
	AppSource string
	TaggedAt  time.Time

	DataloaderKey string
}

type HealthIndicatorEntities []HealthIndicatorEntity

func (HealthIndicatorEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_HealthIndicator,
		NodeLabel_HealthIndicator + "_" + tenant,
	}
}
