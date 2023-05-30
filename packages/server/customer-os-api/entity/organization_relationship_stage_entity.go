package entity

import (
	"time"
)

var (
	DefaultOrganizationRelationshipStageNames = []string{"Target", "Lead", "Prospect", "Trial", "Lost", "Live", "Former"}
)

type OrganizationRelationshipStageEntity struct {
	Id        string
	Name      string
	CreatedAt time.Time
}

type OrganizationRelationshipStageEntities []OrganizationRelationshipStageEntity
