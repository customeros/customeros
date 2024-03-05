package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertIssueCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      model.IssueDataFields
	Source          cmnmod.Source
	ExternalSystem  cmnmod.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertIssueCommand(issueId, tenant, loggedInUserId string, dataFields model.IssueDataFields, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, createdAt, updatedAt *time.Time) *UpsertIssueCommand {
	return &UpsertIssueCommand{
		BaseCommand:    eventstore.NewBaseCommand(issueId, tenant, loggedInUserId),
		DataFields:     dataFields,
		Source:         source,
		ExternalSystem: externalSystem,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

type AddUserAssigneeCommand struct {
	eventstore.BaseCommand
	AssigneeId string `json:"assigneeId" validate:"required"`
	At         *time.Time
	AppSource  string
}

func NewAddUserAssigneeCommand(issueId, tenant, loggedInUserId, userId, appSource string, at *time.Time) *AddUserAssigneeCommand {
	return &AddUserAssigneeCommand{
		BaseCommand: eventstore.NewBaseCommand(issueId, tenant, loggedInUserId),
		AssigneeId:  userId,
		At:          at,
		AppSource:   appSource,
	}
}

type RemoveUserAssigneeCommand struct {
	eventstore.BaseCommand
	AssigneeId string `json:"assigneeId" validate:"required"`
	At         *time.Time
	AppSource  string
}

func NewRemoveUserAssigneeCommand(issueId, tenant, loggedInUserId, userId, appSource string, at *time.Time) *RemoveUserAssigneeCommand {
	return &RemoveUserAssigneeCommand{
		BaseCommand: eventstore.NewBaseCommand(issueId, tenant, loggedInUserId),
		AssigneeId:  userId,
		At:          at,
		AppSource:   appSource,
	}
}

type AddUserFollowerCommand struct {
	eventstore.BaseCommand
	FollowerId string `json:"followerId" validate:"required"`
	At         *time.Time
	AppSource  string
}

func NewAddUserFollowerCommand(issueId, tenant, loggedInUserId, userId, appSource string, at *time.Time) *AddUserFollowerCommand {
	return &AddUserFollowerCommand{
		BaseCommand: eventstore.NewBaseCommand(issueId, tenant, loggedInUserId),
		FollowerId:  userId,
		At:          at,
		AppSource:   appSource,
	}
}

type RemoveUserFollowerCommand struct {
	eventstore.BaseCommand
	FollowerId string `json:"followerId" validate:"required"`
	At         *time.Time
	AppSource  string
}

func NewRemoveUserFollowerCommand(issueId, tenant, loggedInUserId, userId, appSource string, at *time.Time) *RemoveUserFollowerCommand {
	return &RemoveUserFollowerCommand{
		BaseCommand: eventstore.NewBaseCommand(issueId, tenant, loggedInUserId),
		FollowerId:  userId,
		At:          at,
		AppSource:   appSource,
	}
}
