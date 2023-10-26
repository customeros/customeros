package domain

import (
	commentcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/command_handler"
	contactcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command_handler"
	emailcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command_handler"
	iecmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command_handler"
	issuecmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/command_handler"
	jobrolecmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	locationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command_handler"
	logentrycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command_handler"
	organizationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	phonenumbercmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command_handler"
	usercmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command_handler"
)

type CommandHandlers struct {
	Contact          *contactcmdhandler.ContactCommandHandlers
	Organization     *organizationcmdhandler.OrganizationCommandHandlers
	PhoneNumber      *phonenumbercmdhandler.PhoneNumberCommandHandlers
	Email            *emailcmdhandler.EmailCommandHandlers
	User             *usercmdhandler.UserCommandHandlers
	Location         *locationcmdhandler.LocationCommandHandlers
	JobRole          *jobrolecmdhandler.JobRoleCommandHandlers
	InteractionEvent *iecmdhandler.InteractionEventCommandHandlers
	LogEntry         *logentrycmdhandler.LogEntryCommandHandlers
	Issue            *issuecmdhandler.IssueCommandHandlers
	Comment          *commentcmdhandler.CommentCommandHandlers
}
