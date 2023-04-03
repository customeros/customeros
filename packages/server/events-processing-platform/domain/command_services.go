package domain

import (
	contactService "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
)

type Commands struct {
	// FIXME alexb replace with Commands, no need for intermediary service
	ContactCommandsService *contactService.ContactCommandsService
	PhoneNumberCommands    *commands.PhoneNumberCommands
}
