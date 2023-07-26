package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain"
	contactCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/commands"
	emailCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
	jobRoleCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	locationCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/commands"
	organizationCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/commands"
	phoneNumberCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
	userCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

func CreateCommands(log logger.Logger, cfg *config.Config, aggregateStore eventstore.AggregateStore) *domain.Commands {
	return &domain.Commands{
		ContactCommands:      contactCommands.NewContactCommands(log, cfg, aggregateStore),
		OrganizationCommands: organizationCommands.NewOrganizationCommands(log, cfg, aggregateStore),
		PhoneNumberCommands:  phoneNumberCommands.NewPhoneNumberCommands(log, cfg, aggregateStore),
		LocationCommands:     locationCommands.NewLocationCommands(log, cfg, aggregateStore),
		EmailCommands:        emailCommands.NewEmailCommands(log, cfg, aggregateStore),
		UserCommands:         userCommands.NewUserCommands(log, cfg, aggregateStore),
		JobRoleCommands:      jobRoleCommands.NewJobRoleCommands(log, cfg, aggregateStore),
	}
}
