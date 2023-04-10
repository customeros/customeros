package domain

import (
	contact_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/commands"
	email_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
	organization_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/commands"
	phone_number_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
	user_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/commands"
)

type Commands struct {
	ContactCommands      *contact_commands.ContactCommands
	OrganizationCommands *organization_commands.OrganizationCommands
	PhoneNumberCommands  *phone_number_commands.PhoneNumberCommands
	EmailCommands        *email_commands.EmailCommands
	UserCommands         *user_commands.UserCommands
}
