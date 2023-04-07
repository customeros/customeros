package domain

import (
	contact_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/commands"
	email_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
	phone_number_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
)

type Commands struct {
	ContactCommands     *contact_commands.ContactCommands
	PhoneNumberCommands *phone_number_commands.PhoneNumberCommands
	EmailCommands       *email_commands.EmailCommands
}
