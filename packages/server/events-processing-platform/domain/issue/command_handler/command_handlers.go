package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type CommandHandlers struct {
	UpsertIssue        UpsertIssueCommandHandler
	AddUserAssignee    AddUserAssigneeCommandHandler
	RemoveUserAssignee RemoveUserAssigneeCommandHandler
	AddUserFollower    AddUserFollowerCommandHandler
	RemoveUserFollower RemoveUserFollowerCommandHandler
}

func NewCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		UpsertIssue:        NewUpsertIssueCommandHandler(log, es),
		AddUserAssignee:    NewAddUserAssigneeCommandHandler(log, es),
		RemoveUserAssignee: NewRemoveUserAssigneeCommandHandler(log, es),
		AddUserFollower:    NewAddUserFollowerCommandHandler(log, es),
		RemoveUserFollower: NewRemoveUserFollowerCommandHandler(log, es),
	}
}
