package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type InteractionEventCommandHandlers struct {
	UpsertInteractionEvent UpsertInteractionEventCommandHandler
	RequestSummary         RequestSummaryCommandHandler
	ReplaceSummary         ReplaceSummaryCommandHandler
	RequestActionItems     RequestActionItemsCommandHandler
	ReplaceActionItems     ReplaceActionItemsCommandHandler
}

func NewInteractionEventCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *InteractionEventCommandHandlers {
	return &InteractionEventCommandHandlers{
		UpsertInteractionEvent: NewUpsertInteractionEventCommandHandler(log, es),
		RequestSummary:         NewRequestSummaryCommandHandler(log, es),
		ReplaceSummary:         NewReplaceSummaryCommandHandler(log, es),
		RequestActionItems:     NewRequestActionItemsCommandHandler(log, es),
		ReplaceActionItems:     NewReplaceActionItemsCommandHandler(log, es),
	}
}
