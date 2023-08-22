package entity

import (
	"time"
)

type OrganizationRelationshipWithStage struct {
	Relationship  OrganizationRelationship
	Stage         *OrganizationRelationshipStageEntity
	DataloaderKey string
}

type OrganizationRelationshipsWithStages []OrganizationRelationshipWithStage

type OrganizationRelationshipStageEntity struct {
	Id        string
	Name      string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	Order     int64  `neo4jDb:"property:order;lookupName:ORDER;supportCaseSensitive:false"`
	CreatedAt time.Time
}

type OrganizationRelationshipStageEntities []OrganizationRelationshipStageEntity
