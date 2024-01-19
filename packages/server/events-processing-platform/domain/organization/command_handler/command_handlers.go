package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventbuffer"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

// CommandHandlers acts as a container for all command handlers.
type CommandHandlers struct {
	UpsertOrganization           UpsertOrganizationCommandHandler
	UpdateOrganization           UpdateOrganizationCommandHandler
	LinkPhoneNumberCommand       LinkPhoneNumberCommandHandler
	LinkEmailCommand             LinkEmailCommandHandler
	LinkLocationCommand          LinkLocationCommandHandler
	LinkDomainCommand            LinkDomainCommandHandler
	AddSocialCommand             AddSocialCommandHandler
	HideOrganizationCommand      HideOrganizationCommandHandler
	ShowOrganizationCommand      ShowOrganizationCommandHandler
	RefreshLastTouchpointCommand RefreshLastTouchpointCommandHandler
	UpsertCustomFieldCommand     UpsertCustomFieldCommandHandler
	AddParentCommand             AddParentCommandHandler
	RemoveParentCommand          RemoveParentCommandHandler
	RefreshArr                   RefreshArrCommandHandler
	RefreshRenewalSummary        RefreshRenewalSummaryCommandHandler
	WebScrapeOrganization        WebScrapeOrganizationCommandHandler
	UpdateOnboardingStatus       UpdateOnboardingStatusCommandHandler
	UpdateOrganizationOwner      UpdateOrganizationOwnerCommandHandler
}

func NewCommandHandlers(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore, repositories *repository.Repositories, eventBufferWatcher *eventbuffer.EventBufferWatcher) *CommandHandlers {
	return &CommandHandlers{
		UpsertOrganization:           NewUpsertOrganizationCommandHandler(log, es),
		UpdateOrganization:           NewUpdateOrganizationCommandHandler(log, es, cfg.Utils),
		LinkPhoneNumberCommand:       NewLinkPhoneNumberCommandHandler(log, es),
		LinkEmailCommand:             NewLinkEmailCommandHandler(log, es),
		LinkLocationCommand:          NewLinkLocationCommandHandler(log, es),
		LinkDomainCommand:            NewLinkDomainCommandHandler(log, es, cfg.Utils),
		AddSocialCommand:             NewAddSocialCommandHandler(log, es, cfg.Utils),
		HideOrganizationCommand:      NewHideOrganizationCommandHandler(log, es),
		ShowOrganizationCommand:      NewShowOrganizationCommandHandler(log, es),
		RefreshLastTouchpointCommand: NewRefreshLastTouchpointCommandHandler(log, es, cfg.Utils),
		UpsertCustomFieldCommand:     NewUpsertCustomFieldCommandHandler(log, es),
		AddParentCommand:             NewAddParentCommandHandler(log, es),
		RemoveParentCommand:          NewRemoveParentCommandHandler(log, es),
		RefreshArr:                   NewRefreshArrCommandHandler(log, es, cfg.Utils),
		RefreshRenewalSummary:        NewRefreshRenewalSummaryCommandHandler(log, es, cfg.Utils),
		WebScrapeOrganization:        NewWebScrapeOrganizationCommandHandler(log, es, cfg.Utils),
		UpdateOnboardingStatus:       NewUpdateOnboardingStatusCommandHandler(log, es, cfg.Utils),
		UpdateOrganizationOwner:      NewUpdateOrganizationOwnerCommandHandler(log, es, cfg.Utils, eventBufferWatcher),
	}
}
