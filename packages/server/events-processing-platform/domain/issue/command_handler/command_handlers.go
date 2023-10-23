package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type IssueCommandHandlers struct {
	UpsertIssue        UpsertIssueCommandHandler
	AddUserAssignee    AddUserAssigneeCommandHandler
	RemoveUserAssignee RemoveUserAssigneeCommandHandler
	AddUserFollower    AddUserFollowerCommandHandler
	RemoveUserFollower RemoveUserFollowerCommandHandler
}

func NewIssueCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *IssueCommandHandlers {
	return &IssueCommandHandlers{
		UpsertIssue:        NewUpsertIssueCommandHandler(log, es),
		AddUserAssignee:    NewAddUserAssigneeCommandHandler(log, es),
		RemoveUserAssignee: NewRemoveUserAssigneeCommandHandler(log, es),
		AddUserFollower:    NewAddUserFollowerCommandHandler(log, es),
		RemoveUserFollower: NewRemoveUserFollowerCommandHandler(log, es),
	}
}
