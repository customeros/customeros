package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type PhoneNumberCommands struct {
	UpsertPhoneNumber UpsertPhoneNumberCommandHandler
	CreatePhoneNumber CreatePhoneNumberCommandHandler
	// alexbalexb
	//FailPhoneNumberValidation FailPhoneNumberValidationCommandHandler
	//PhoneNumberValidated      PhoneNumberValidatedCommandHandler
}

func NewPhoneNumberCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *PhoneNumberCommands {
	return &PhoneNumberCommands{
		CreatePhoneNumber: NewCreatePhoneNumberCommandHandler(log, cfg, es),
		UpsertPhoneNumber: NewUpsertPhoneNumberHandler(log, cfg, es),
	}
}
