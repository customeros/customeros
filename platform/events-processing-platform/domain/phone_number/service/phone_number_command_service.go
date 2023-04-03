package service

import (
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/phone_number/commands"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/logger"
)

type PhoneNumberCommandsService struct {
	Commands *commands.PhoneNumberCommands
}

func NewPhoneNumberCommandsService(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *PhoneNumberCommandsService {
	createPhoneNumberHandler := commands.NewCreatePhoneNumberHandler(log, cfg, es)

	contactCommands := commands.NewPhoneNumberCommands(
		createPhoneNumberHandler,
	)

	return &PhoneNumberCommandsService{Commands: contactCommands}
}
