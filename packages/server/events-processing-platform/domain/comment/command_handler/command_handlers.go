package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type CommentCommandHandlers struct {
	Upsert UpsertCommentCommandHandler
}

func NewCommentCommandHandlers(log logger.Logger, es eventstore.AggregateStore) *CommentCommandHandlers {
	return &CommentCommandHandlers{
		Upsert: NewUpsertCommentCommandHandler(log, es),
	}
}
