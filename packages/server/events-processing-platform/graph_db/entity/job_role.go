package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type JobRoleEntity struct {
	Id            string
	JobTitle      string
	Primary       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StartedAt     *time.Time
	EndedAt       *time.Time
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string
	Description   *string
	Company       *string

	InteractionEventParticipantDetails InteractionEventParticipantDetails

	DataloaderKey string
}

type JobRoleEntities []JobRoleEntity

func (jobRole JobRoleEntity) ToString() string {
	return fmt.Sprintf("id: %s\njob title: %s", jobRole.Id, jobRole.JobTitle)
}

func (JobRoleEntity) IsInteractionEventParticipant() {}

func (JobRoleEntity) ParticipantLabel() string {
	return NodeLabel_JobRole
}

func (jobRole JobRoleEntity) GetDataloaderKey() string {
	return jobRole.DataloaderKey
}

func (jobRole JobRoleEntity) Labels(tenant string) []string {
	return []string{"JobRole", "JobRole_" + tenant}
}
