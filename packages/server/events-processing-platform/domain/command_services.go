package domain

import (
	contact_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/commands"
	email_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
	interaction_event_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/commands"
	job_role_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	location_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/commands"
	log_entry_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command_handler"
	orgcmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	phone_number_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
	user_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/commands"
)

type Commands struct {
	ContactCommands          *contact_commands.ContactCommands
	OrganizationCommands     *orgcmdhnd.OrganizationCommands
	PhoneNumberCommands      *phone_number_commands.PhoneNumberCommands
	EmailCommands            *email_commands.EmailCommands
	UserCommands             *user_commands.UserCommands
	LocationCommands         *location_commands.LocationCommands
	JobRoleCommands          *job_role_commands.JobRoleCommands
	InteractionEventCommands *interaction_event_commands.InteractionEventCommands
	LogEntryCommands         *log_entry_commands.LogEntryCommands
}
