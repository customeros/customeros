package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type EmailEntity struct {
	Id            string
	Email         string `neo4jDb:"property:email;lookupName:EMAIL;supportCaseSensitive:true"`
	RawEmail      string `neo4jDb:"property:rawEmail;lookupName:RAW_EMAIL;supportCaseSensitive:true"`
	Label         string
	Primary       bool
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Validated      *bool
	IsReachable    *string
	IsValidSyntax  *bool
	CanConnectSMTP *bool
	AcceptsMail    *bool
	HasFullInbox   *bool
	IsCatchAll     *bool
	IsDeliverable  *bool
	IsDisabled     *bool
	Error          *string
}

type EmailEntities []EmailEntity
