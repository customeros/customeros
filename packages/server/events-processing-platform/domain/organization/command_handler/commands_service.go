package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

type OrganizationCommandHandlers struct {
	UpsertOrganization             UpsertOrganizationCommandHandler
	UpdateOrganization             UpdateOrganizationCommandHandler
	LinkPhoneNumberCommand         LinkPhoneNumberCommandHandler
	LinkEmailCommand               LinkEmailCommandHandler
	LinkLocationCommand            LinkLocationCommandHandler
	LinkDomainCommand              LinkDomainCommandHandler
	AddSocialCommand               AddSocialCommandHandler
	UpdateRenewalLikelihoodCommand UpdateRenewalLikelihoodCommandHandler
	UpdateRenewalForecastCommand   UpdateRenewalForecastCommandHandler
	UpdateBillingDetailsCommand    UpdateBillingDetailsCommandHandler
	RequestRenewalForecastCommand  RequestRenewalForecastCommandHandler
	RequestNextCycleDateCommand    RequestNextCycleDateCommandHandler
	HideOrganizationCommand        HideOrganizationCommandHandler
	ShowOrganizationCommand        ShowOrganizationCommandHandler
	RefreshLastTouchpointCommand   RefreshLastTouchpointCommandHandler
	UpsertCustomFieldCommand       UpsertCustomFieldCommandHandler
	AddParentCommand               AddParentCommandHandler
	RemoveParentCommand            RemoveParentCommandHandler
}

func NewOrganizationCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore, repositories *repository.Repositories) *OrganizationCommandHandlers {
	return &OrganizationCommandHandlers{
		UpsertOrganization:             NewUpsertOrganizationCommandHandler(log, es),
		UpdateOrganization:             NewUpdateOrganizationCommandHandler(log, es),
		LinkPhoneNumberCommand:         NewLinkPhoneNumberCommandHandler(log, es),
		LinkEmailCommand:               NewLinkEmailCommandHandler(log, es),
		LinkLocationCommand:            NewLinkLocationCommandHandler(log, es),
		LinkDomainCommand:              NewLinkDomainCommandHandler(log, es),
		AddSocialCommand:               NewAddSocialCommandHandler(log, es),
		UpdateRenewalLikelihoodCommand: NewUpdateRenewalLikelihoodCommandHandler(log, es, repositories),
		UpdateRenewalForecastCommand:   NewUpdateRenewalForecastCommandHandler(log, es, repositories),
		UpdateBillingDetailsCommand:    NewUpdateBillingDetailsCommandHandler(log, es, repositories),
		RequestRenewalForecastCommand:  NewRequestRenewalForecastCommandHandler(log, es),
		RequestNextCycleDateCommand:    NewRequestNextCycleDateCommandHandler(log, es),
		HideOrganizationCommand:        NewHideOrganizationCommandHandler(log, es),
		ShowOrganizationCommand:        NewShowOrganizationCommandHandler(log, es),
		RefreshLastTouchpointCommand:   NewRefreshLastTouchpointCommandHandler(log, es),
		UpsertCustomFieldCommand:       NewUpsertCustomFieldCommandHandler(log, es),
		AddParentCommand:               NewAddParentCommandHandler(log, es),
		RemoveParentCommand:            NewRemoveParentCommandHandler(log, es),
	}
}
