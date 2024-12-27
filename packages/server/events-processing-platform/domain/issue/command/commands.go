package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"time"
)

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
