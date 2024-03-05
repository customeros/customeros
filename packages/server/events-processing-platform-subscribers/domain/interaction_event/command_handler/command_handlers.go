package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type CommandHandlers struct {
	UpsertInteractionEvent UpsertInteractionEventCommandHandler
	RequestSummary         RequestSummaryCommandHandler
	ReplaceSummary         ReplaceSummaryCommandHandler
	RequestActionItems     RequestActionItemsCommandHandler
	ReplaceActionItems     ReplaceActionItemsCommandHandler
}

func NewCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		UpsertInteractionEvent: NewUpsertInteractionEventCommandHandler(log, es),
		RequestSummary:         NewRequestSummaryCommandHandler(log, es),
		ReplaceSummary:         NewReplaceSummaryCommandHandler(log, es),
		RequestActionItems:     NewRequestActionItemsCommandHandler(log, es),
		ReplaceActionItems:     NewReplaceActionItemsCommandHandler(log, es),
	}
}
