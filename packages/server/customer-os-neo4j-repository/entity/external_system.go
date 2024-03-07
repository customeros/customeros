package entity

import (
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type ExternalSystemEntity struct {
	ExternalSystemId neo4jenum.ExternalSystemId
	Relationship     struct {
		ExternalId     string
		SyncDate       *time.Time
		ExternalUrl    *string
		ExternalSource *string
	}
	DataloaderKey string
}

type ExternalSystemEntities []ExternalSystemEntity
