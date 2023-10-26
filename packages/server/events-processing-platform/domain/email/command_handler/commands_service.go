package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type EmailCommandHandlers struct {
	Upsert              UpsertEmailCommandHandler
	FailEmailValidation FailEmailValidationCommandHandler
	EmailValidated      EmailValidatedCommandHandler
}

func NewEmailCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *EmailCommandHandlers {
	return &EmailCommandHandlers{
		Upsert:              NewUpsertEmailHandler(log, cfg, es),
		FailEmailValidation: NewFailEmailValidationCommandHandler(log, cfg, es),
		EmailValidated:      NewEmailValidatedCommandHandler(log, cfg, es),
	}
}
