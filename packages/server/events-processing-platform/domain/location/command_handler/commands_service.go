package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type LocationCommands struct {
	UpsertLocation           UpsertLocationCommandHandler
	FailedLocationValidation FailedLocationValidationCommandHandler
	SkipLocationValidation   SkippedLocationValidationCommandHandler
	LocationValidated        LocationValidatedCommandHandler
}

func NewLocationCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *LocationCommands {
	return &LocationCommands{
		UpsertLocation:           NewUpsertLocationHandler(log, cfg, es),
		FailedLocationValidation: NewFailedLocationValidationCommandHandler(log, cfg, es),
		SkipLocationValidation:   NewSkippedLocationValidationCommandHandler(log, cfg, es),
		LocationValidated:        NewLocationValidatedCommandHandler(log, cfg, es),
	}
}
