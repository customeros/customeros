package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type CommandHandlers struct {
	UpsertLocation           UpsertLocationCommandHandler
	FailedLocationValidation FailedLocationValidationCommandHandler
	SkipLocationValidation   SkippedLocationValidationCommandHandler
	LocationValidated        LocationValidatedCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		UpsertLocation:           NewUpsertLocationHandler(log, cfg, es),
		FailedLocationValidation: NewFailedLocationValidationCommandHandler(log, cfg, es),
		SkipLocationValidation:   NewSkippedLocationValidationCommandHandler(log, cfg, es),
		LocationValidated:        NewLocationValidatedCommandHandler(log, cfg, es),
	}
}
