package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type CommandHandlers struct {
	UpsertPhoneNumber           UpsertPhoneNumberCommandHandler
	FailedPhoneNumberValidation FailedPhoneNumberValidationCommandHandler
	SkipPhoneNumberValidation   SkippedPhoneNumberValidationCommandHandler
	PhoneNumberValidated        PhoneNumberValidatedCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *CommandHandlers {
	return &CommandHandlers{
		UpsertPhoneNumber:           NewUpsertPhoneNumberHandler(log, cfg, es),
		FailedPhoneNumberValidation: NewFailedPhoneNumberValidationCommandHandler(log, cfg, es),
		SkipPhoneNumberValidation:   NewSkippedPhoneNumberValidationCommandHandler(log, cfg, es),
		PhoneNumberValidated:        NewPhoneNumberValidatedCommandHandler(log, cfg, es),
	}
}
