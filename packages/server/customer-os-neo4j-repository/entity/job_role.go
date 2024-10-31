package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type JobRoleEntity struct {
	DataLoaderKey
	Id          string
	JobTitle    string
	Primary     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	StartedAt   *time.Time
	EndedAt     *time.Time
	Source      DataSource
	AppSource   string
	Description *string
	Company     *string

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails
}

type JobRoleEntities []JobRoleEntity

func (e JobRoleEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (JobRoleEntity) IsInteractionEventParticipant() {}

func (JobRoleEntity) IsInteractionSessionParticipant() {}

func (JobRoleEntity) EntityLabel() string {
	return model.NodeLabelJobRole
}

func (JobRoleEntity) ParticipantLabel() string {
	return model.NodeLabelJobRole
}
