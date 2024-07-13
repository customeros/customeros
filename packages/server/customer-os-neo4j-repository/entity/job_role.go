package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type JobRoleEntity struct {
	DataLoaderKey
	Id            string
	JobTitle      string
	Primary       bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	StartedAt     *time.Time
	EndedAt       *time.Time
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	Description   *string
	Company       *string

	InteractionEventParticipantDetails InteractionEventParticipantDetails
}

type JobRoleEntities []JobRoleEntity

func (JobRoleEntity) IsInteractionEventParticipant() {}

func (JobRoleEntity) ParticipantLabel() string {
	return model.NodeLabelJobRole
}
