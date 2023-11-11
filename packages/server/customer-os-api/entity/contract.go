package entity

import (
	"time"
)

type ContractEntity struct {
	ID                   string
	Name                 string     `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	CreatedAt            *time.Time `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt            *time.Time `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	ServiceStartedAt     *time.Time `neo4jDb:"property:serviceStartedAt;lookupName:SERVICE_STARTED_AT;supportCaseSensitive:false"`
	SignedAt             *time.Time `neo4jDb:"property:SignedAt;lookupName:SIGNED_AT;supportCaseSensitive:false"`
	EndedAt              *time.Time `neo4jDb:"property:EndedAt;lookupName:ENDED_AT;supportCaseSensitive:false"`
	ContractRenewalCycle ContractRenewalCycle
	ContractStatus       ContractStatus
	Source               DataSource
	AppSource            string
	ContractUrl          string
	DataloaderKey        string
	CreatedByUsedId      string
}
