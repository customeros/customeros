package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type PhoneNumberEntity struct {
	Id             string
	E164           string
	RawPhoneNumber string
	Source         neo4jentity.DataSource
	SourceOfTruth  neo4jentity.DataSource
	AppSource      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type PhoneNumberEntities []PhoneNumberEntity
