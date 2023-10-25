package domain

import (
	contactcmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command_handler"
	emailcmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command_handler"
	iecmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command_handler"
	issuecmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/command_handler"
	jobrolecmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	locationcmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command_handler"
	logentrycmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command_handler"
	orgcmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	phonecmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command_handler"
	usercmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command_handler"
)

type Commands struct {
	ContactCommands          *contactcmd.ContactCommands
	OrganizationCommands     *orgcmd.OrganizationCommands
	PhoneNumberCommands      *phonecmd.PhoneNumberCommands
	EmailCommands            *emailcmd.EmailCommands
	UserCommands             *usercmd.UserCommands
	LocationCommands         *locationcmd.LocationCommands
	JobRoleCommands          *jobrolecmd.JobRoleCommands
	InteractionEventCommands *iecmd.InteractionEventCommandHandlers
	LogEntryCommands         *logentrycmd.LogEntryCommandHandlers
	IssueCommands            *issuecmd.IssueCommandHandlers
}
