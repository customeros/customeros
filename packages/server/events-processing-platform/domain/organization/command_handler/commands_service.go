package command_handler

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

type OrganizationCommands struct {
	UpsertOrganization             UpsertOrganizationCommandHandler
	UpdateOrganization             UpdateOrganizationCommandHandler
	LinkPhoneNumberCommand         LinkPhoneNumberCommandHandler
	LinkEmailCommand               LinkEmailCommandHandler
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
}

func NewOrganizationCommands(log logger.Logger, cfg *config.Config, es eventstore.AggregateStore, repositories *repository.Repositories) *OrganizationCommands {
	return &OrganizationCommands{
		UpsertOrganization:             NewUpsertOrganizationCommandHandler(log, cfg, es),
		UpdateOrganization:             NewUpdateOrganizationCommandHandler(log, cfg, es),
		LinkPhoneNumberCommand:         NewLinkPhoneNumberCommandHandler(log, cfg, es),
		LinkEmailCommand:               NewLinkEmailCommandHandler(log, cfg, es),
		LinkDomainCommand:              NewLinkDomainCommandHandler(log, cfg, es),
		AddSocialCommand:               NewAddSocialCommandHandler(log, cfg, es),
		UpdateRenewalLikelihoodCommand: NewUpdateRenewalLikelihoodCommandHandler(log, cfg, es, repositories),
		UpdateRenewalForecastCommand:   NewUpdateRenewalForecastCommandHandler(log, cfg, es, repositories),
		UpdateBillingDetailsCommand:    NewUpdateBillingDetailsCommandHandler(log, cfg, es, repositories),
		RequestRenewalForecastCommand:  NewRequestRenewalForecastCommandHandler(log, es),
		RequestNextCycleDateCommand:    NewRequestNextCycleDateCommandHandler(log, es),
		HideOrganizationCommand:        NewHideOrganizationCommandHandler(log, es),
		ShowOrganizationCommand:        NewShowOrganizationCommandHandler(log, es),
		RefreshLastTouchpointCommand:   NewRefreshLastTouchpointCommandHandler(log, es),
		UpsertCustomFieldCommand:       NewUpsertCustomFieldCommandHandler(log, es),
	}
}
