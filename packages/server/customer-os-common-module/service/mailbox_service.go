package service

import (
	"context"
	"gorm.io/gorm"
	"strings"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/coserrors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type CreateMailboxRequest struct {
	Domain          string
	Username        string
	Password        string
	LinkedUserEmail string
	WebmailEnabled  bool
	ForwardingTo    []string

	IgnoreDomainOwnership bool
}

type MailboxService interface {
	CreateMailbox(ctx context.Context, tx *gorm.DB, request CreateMailboxRequest) error
	IsDomainAvailable(ctx context.Context, domain string) (ok, available bool)
	RecommendOutboundDomains(ctx context.Context, domainRoot string, count int) []string
	ReputationScore(ctx context.Context, domain, tenant string) (int, error)
}

type mailboxService struct {
	log      logger.Logger
	services *Services
}

const TEST_MAILBOX_DOMAIN = "testcustomeros.com"

func NewMailboxService(log logger.Logger, services *Services) MailboxService {
	return &mailboxService{
		log:      log,
		services: services,
	}
}

func (s *mailboxService) CreateMailbox(ctx context.Context, tx *gorm.DB, request CreateMailboxRequest) error {
	span, ctx := s.initializeTracing(ctx, "MailboxService.CreateMailbox")
	span.LogFields(
		log.String("linkedUserEmail", request.LinkedUserEmail),
		log.String("domain", request.Domain),
		log.String("username", request.Username),
		log.Bool("webmailEnabled", request.WebmailEnabled),
		log.Object("forwardingTo", request.ForwardingTo),
	)
	defer span.Finish()

	if !request.IgnoreDomainOwnership {
		if err := s.validateRequest(ctx, span, request.Domain); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "cannot vaildate MailboxRequest"))
			return err
		}
	}

	mailboxEmail := request.Username + "@" + request.Domain

	// Find linked user if email provided
	userId, linkedUserFound, err := s.findLinkedUser(ctx, span, request.LinkedUserEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "cannot find linked user"))
		return err
	}

	// Verify mailbox doesn't exist
	if err := s.verifyMailboxNotExists(ctx, span, mailboxEmail); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to verify mailbox does not exist"))
		return err
	}

	// Save mailbox
	if err := s.createMailbox(ctx, span, tx, request, mailboxEmail, userId); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to save mailbox settings"))
		return err
	}

	// Create email node
	if err := s.createEmailNode(ctx, span, mailboxEmail, linkedUserFound, userId); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create email node"))
		return err
	}

	return nil
}

func (s *mailboxService) initializeTracing(ctx context.Context, methodName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, methodName)
	tracing.SetDefaultServiceSpanTags(ctx, span)
	return span, ctx
}

func (s *mailboxService) validateRequest(ctx context.Context, span opentracing.Span, domain string) error {
	if err := common.ValidateTenant(ctx); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error validating tenant"))
		return err
	}

	tenant := common.GetTenantFromContext(ctx)
	if domain != TEST_MAILBOX_DOMAIN {
		domainBelongsToTenant, err := s.services.PostgresRepositories.MailStackDomainRepository.CheckDomainOwnership(ctx, tenant, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error checking domain"))
			return errors.Wrap(err, "Error checking domain")
		}
		if !domainBelongsToTenant {
			tracing.TraceErr(span, errors.Wrap(coserrors.ErrDomainNotFound, "domain does not belong to tenant"))
			return coserrors.ErrDomainNotFound
		}
	}
	return nil
}

func (s *mailboxService) findLinkedUser(ctx context.Context, span opentracing.Span, linkedUserEmail string) (string, bool, error) {
	if linkedUserEmail == "" {
		return "", false, nil
	}

	tenant := common.GetTenantFromContext(ctx)
	userDbNode, err := s.services.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, linkedUserEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get user by email"))
		return "", false, err
	}
	if userDbNode == nil {
		return "", false, nil
	}

	userEntity := neo4jmapper.MapDbNodeToUserEntity(userDbNode)
	return userEntity.Id, true, nil
}

func (s *mailboxService) verifyMailboxNotExists(ctx context.Context, span opentracing.Span, mailboxEmail string) error {
	mailboxRecord, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, mailboxEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error checking mailbox"))
		return err
	}
	if mailboxRecord != nil {
		tracing.TraceErr(span, errors.Wrap(coserrors.ErrMailboxExists, "Mailbox exists"))
		return coserrors.ErrMailboxExists
	}
	return nil
}

func (s *mailboxService) createMailbox(ctx context.Context, span opentracing.Span, tx *gorm.DB, request CreateMailboxRequest, mailboxEmail string, userId string) error {
	tenant := common.GetTenantFromContext(ctx)
	tenantSettingsMailbox := entity.TenantSettingsMailbox{
		Tenant:          tenant,
		Domain:          request.Domain,
		MailboxUsername: mailboxEmail,
		MailboxPassword: request.Password,
		Username:        request.Username,
		UserId:          userId,
		ForwardingTo:    strings.Join(request.ForwardingTo, ","),
		WebmailEnabled:  request.WebmailEnabled,
		Status:          entity.MailboxStatusPendingProvisioning,
	}
	err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.Merge(ctx, tx, &tenantSettingsMailbox)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error saving mailbox"))
		return err
	}
	return nil
}

func (s *mailboxService) createEmailNode(ctx context.Context, span opentracing.Span, mailboxEmail string, linkedUserFound bool, userId string) error {
	tenant := common.GetTenantFromContext(ctx)
	emailFields := EmailFields{
		Email:     mailboxEmail,
		Source:    neo4jentity.DataSourceOpenline,
		AppSource: common.GetAppSourceFromContext(ctx),
	}

	var linkWith *LinkWith
	if linkedUserFound {
		linkWith = &LinkWith{
			Type: model.USER,
			Id:   userId,
		}
	}

	_, err := s.services.EmailService.Merge(ctx, nil, tenant, emailFields, linkWith)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating email node for mailbox"))
		return err
	}
	return nil
}
