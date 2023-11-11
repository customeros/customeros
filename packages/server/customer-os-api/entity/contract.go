package entity

import (
	"time"
)

type ContractEntity struct {
	ID                                 string
	CustomerOsId                       string     `neo4jDb:"property:customerOsId;lookupName:CUSTOMER_OS_ID;supportCaseSensitive:false"`
	ReferenceId                        string     `neo4jDb:"property:referenceId;lookupName:REFERENCE_ID;supportCaseSensitive:true"`
	Name                               string     `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	CreatedAt                          time.Time  `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt                          time.Time  `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	ServiceStartedAt                   *time.Time `neo4jDb:"property:serviceStartedAt;lookupName:SERVICE_STARTED_AT;supportCaseSensitive:false"`
	ServiceStartedId                   *string    `neo4jDb:"property:serviceStartedId;lookupName:SERVICE_STARTED_ID;supportCaseSensitive:false"`
	ContractRenewalCycle               ContractRenewalCycle
	ContractStatus                     ContractStatus
	Source                             DataSource
	AppSource                          string
	InteractionEventParticipantDetails InteractionEventParticipantDetails
	DataloaderKey                      string
}
