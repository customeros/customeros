package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/coserrors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

const TEST_MAILBOX_DOMAIN = "testcustomeros.com"

type MailboxService interface {
	AddMailbox(ctx context.Context, domain, username, password string, forwardingEnabled, webmailEnabled bool, forwardingTo []string) error
}

type mailboxService struct {
	log      logger.Logger
	services *Services
}

func NewMailboxService(log logger.Logger, services *Services) MailboxService {
	return &mailboxService{
		log:      log,
		services: services,
	}
}

func (s *mailboxService) AddMailbox(ctx context.Context, domain, username, password string, forwardingEnabled, webmailEnabled bool, forwardingTo []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailboxService.SaveContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("domain", domain), log.String("username", username), log.Bool("forwardingEnabled", forwardingEnabled), log.Bool("webmailEnabled", webmailEnabled), log.Object("forwardingTo", forwardingTo))

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	mailboxEmail := username + "@" + domain

	// Check domain belongs to tenant, skip verification for test domain
	if domain != TEST_MAILBOX_DOMAIN {
		domainBelongsToTenant, err := s.services.PostgresRepositories.MailStackDomainRepository.CheckDomainOwnership(ctx, tenant, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error checking domain"))
			return err
		}
		if !domainBelongsToTenant {
			return coserrors.ErrDomainNotFound
		}
	}

	// Check mailbox doesn't already exist
	mailboxRecord, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, mailboxEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error checking mailbox"))
		return err
	}
	if mailboxRecord != nil {
		return coserrors.ErrMailboxExists
	}

	err = s.services.OpenSrsService.SetMailbox(ctx, tenant, domain, username, password, forwardingEnabled, forwardingTo, webmailEnabled)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error setting mailbox"))
		return err
	}

	// Save mailbox details in postgres
	err = s.services.PostgresRepositories.TenantSettingsMailboxRepository.Merge(ctx, &entity.TenantSettingsMailbox{
		Domain: domain, MailboxUsername: mailboxEmail, Tenant: tenant, MailboxPassword: password, MinMinutesBetweenEmails: 5, MaxMinutesBetweenEmails: 10,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error saving mailbox"))
		return err
	}

	// create email node in neo4j
	_, err = s.services.EmailService.Merge(ctx, tenant, EmailFields{
		Email:     mailboxEmail,
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: common.GetAppSourceFromContext(ctx),
	}, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating email node for mailbox"))
	}

	return nil

}
