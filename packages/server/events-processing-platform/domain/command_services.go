package domain

import (
	contactCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/commands"
	phoneNumberCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
)

type Commands struct {
	ContactCommands     *contactCommands.ContactCommands
	PhoneNumberCommands *phoneNumberCommands.PhoneNumberCommands
}
