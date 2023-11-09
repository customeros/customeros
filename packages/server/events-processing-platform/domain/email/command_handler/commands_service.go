package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type CommandHandlers struct {
	Upsert              UpsertEmailCommandHandler
	FailEmailValidation FailEmailValidationCommandHandler
	EmailValidated      EmailValidatedCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		Upsert:              NewUpsertEmailHandler(log, cfg, es),
		FailEmailValidation: NewFailEmailValidationCommandHandler(log, cfg, es),
		EmailValidated:      NewEmailValidatedCommandHandler(log, cfg, es),
	}
}
