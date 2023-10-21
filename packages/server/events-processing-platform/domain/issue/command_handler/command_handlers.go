package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type IssueCommandHandlers struct {
	UpsertIssue UpsertIssueCommandHandler
}

func NewIssueCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *IssueCommandHandlers {
	return &IssueCommandHandlers{
		UpsertIssue: NewUpsertIssueCommandHandler(log, cfg, es),
	}
}
