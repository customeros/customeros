package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type HealthIndicatorEntity struct {
	Id        string
	Name      string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	Order     int64  `neo4jDb:"property:order;lookupName:ORDER;supportCaseSensitive:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Source    neo4jentity.DataSource
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
