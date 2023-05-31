package entity

import (
	"time"
)

var (
	DefaultOrganizationRelationshipStageNames = []string{"Target", "Lead", "Prospect", "Trial", "Lost", "Live", "Former"}
)

type OrganizationRelationshipWithStage struct {
	Relationship  OrganizationRelationship
	Stage         *OrganizationRelationshipStageEntity
	DataloaderKey string
}

type OrganizationRelationshipsWithStages []OrganizationRelationshipWithStage

type OrganizationRelationshipStageEntity struct {
	Id        string
	Name      string
	CreatedAt time.Time
}

type OrganizationRelationshipStageEntities []OrganizationRelationshipStageEntity
