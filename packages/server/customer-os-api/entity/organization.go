package entity

import (
	"fmt"
	"time"
)

type OrganizationEntity struct {
	ID                 string
	TenantOrganization bool
	Name               string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	Description        string `neo4jDb:"property:description;lookupName:DESCRIPTION;supportCaseSensitive:true"`
	Website            string `neo4jDb:"property:website;lookupName:WEBSITE;supportCaseSensitive:true"`
	Industry           string `neo4jDb:"property:industry;lookupName:INDUSTRY;supportCaseSensitive:true"`
	SubIndustry        string
	IndustryGroup      string
	TargetAudience     string
	ValueProposition   string
	IsPublic           bool
	Market             string
	LastFundingRound   string
	LastFundingAmount  string
	Employees          int64
	CreatedAt          time.Time  `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt          time.Time  `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	LastTouchpointAt   *time.Time `neo4jDb:"property:lastTouchpointAt;lookupName:LAST_TOUCHPOINT_AT;supportCaseSensitive:false"`
	LastTouchpointId   *string    `neo4jDb:"property:lastTouchpointId;lookupName:LAST_TOUCHPOINT_ID;supportCaseSensitive:false"`
	Source             DataSource
	SourceOfTruth      DataSource
	AppSource          string

	LinkedOrganizationType *string

	InteractionEventParticipantDetails InteractionEventParticipantDetails

	DataloaderKey string
}

func (organization OrganizationEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", organization.ID, organization.Name)
}

func (OrganizationEntity) IsNotedEntity() {}

func (OrganizationEntity) NotedEntityLabel() string {
	return NodeLabel_Organization
}

func (OrganizationEntity) IsInteractionEventParticipant() {}

func (OrganizationEntity) InteractionEventParticipantLabel() string {
	return NodeLabel_Organization
}

func (OrganizationEntity) IsMeetingParticipant() {}

func (OrganizationEntity) MeetingParticipantLabel() string {
	return NodeLabel_Organization
}

func (organization OrganizationEntity) GetDataloaderKey() string {
	return organization.DataloaderKey
}

type OrganizationEntities []OrganizationEntity

func (organization OrganizationEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Organization,
		NodeLabel_Organization + "_" + tenant,
	}
}
