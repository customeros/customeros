package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type PhoneNumberCommands struct {
	CreatePhoneNumber CreatePhoneNumberCommandHandler
	UpsertPhoneNumber UpsertPhoneNumberCommandHandler
}

func NewPhoneNumberCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore) *PhoneNumberCommands {
	return &PhoneNumberCommands{
		CreatePhoneNumber: NewCreatePhoneNumberHandler(log, cfg, es),
		UpsertPhoneNumber: NewUpsertPhoneNumberHandler(log, cfg, es),
	}
}
