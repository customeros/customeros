package service

import (
	"context"
	"fmt"
	"strings"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type RegistrationService interface {
	PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error
	ConfigureTestMailbox(ctx context.Context) error
	CreatePostmarkServer(ctx context.Context) error
}

type registrationService struct {
	services *Services
}

func NewRegistrationService(services *Services) RegistrationService {
	return &registrationService{
		services: services,
	}
}

type testUserSetup struct {
	userId         string
	mailboxAddress string
}

func (s *registrationService) PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error {
	span, ctx := s.initializeTracing(ctx, "PrepareDefaultTenantSetup", map[string]interface{}{
		"loggedInUserEmail": loggedInUserEmail,
	})
	defer span.Finish()

	if err := common.ValidateTenant(ctx); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err := s.ConfigureTestMailbox(ctx); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error configuring test mailbox during tenant onboarding"))
	}

	if err := s.CreatePostmarkServer(ctx); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating postmark server during tenant onboarding"))
	}

	return nil
}

func (s *registrationService) ConfigureTestMailbox(ctx context.Context) error {
	span, ctx := s.initializeTracing(ctx, "ConfigureTestMailbox", nil)
	defer span.Finish()

	if err := common.ValidateTenant(ctx); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	tenant := common.GetTenantFromContext(ctx)
	testUser, err := s.setupTestUser(ctx, span)
	if err != nil {
		return err
	}

	if err := s.setupTestMailbox(ctx, span, tenant, testUser); err != nil {
		return err
	}

	span.LogKV("result.mailboxAddress", testUser.mailboxAddress)
	return nil
}

func (s *registrationService) CreatePostmarkServer(ctx context.Context) error {
	span, ctx := s.initializeTracing(ctx, "CreatePostmarkServer", nil)
	defer span.Finish()

	if err := common.ValidateTenant(ctx); err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if err := s.services.PostmarkService.CreateServer(ctx); err != nil {
		tracing.TraceErr(span, err)
	}

	return nil
}

// Helper functions

func (s *registrationService) initializeTracing(ctx context.Context, operation string, logFields map[string]interface{}) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("RegistrationService.%s", operation))
	tracing.SetDefaultServiceSpanTags(ctx, span)
	for key, value := range logFields {
		span.LogKV(key, value)
	}
	return span, ctx
}

func (s *registrationService) setupTestUser(ctx context.Context, span opentracing.Span) (*testUserSetup, error) {
	existingTestUser, err := s.services.Neo4jRepositories.UserReadRepository.FindTestUser(ctx)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "cannot find test user"))
		return nil, err
	}

	var testUserId string
	if existingTestUser == nil {
		testUserId, err = s.services.UserService.CreateUser(ctx, neo4jentity.UserEntity{
			FirstName: "Test",
			LastName:  "Sender",
			Test:      true,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "cannot create test user"))
			return nil, err
		}
	} else {
		testUserId = mapper.MapDbNodeToUserEntity(existingTestUser).Id
	}

	span.LogKV("result.testUserId", testUserId)
	return &testUserSetup{userId: testUserId}, nil
}

func (s *registrationService) setupTestMailbox(ctx context.Context, span opentracing.Span, tenant string, testUser *testUserSetup) error {
	mailboxAddress := strings.ToLower(fmt.Sprintf("%s@%s", tenant, TEST_MAILBOX_DOMAIN))
	testUser.mailboxAddress = mailboxAddress

	testEmailId, err := s.services.EmailService.Merge(ctx, nil, tenant, EmailFields{
		Email: mailboxAddress,
	}, &LinkWith{
		Type: model.USER,
		Id:   testUser.userId,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to setup test mailbox"))
		return err
	}
	span.LogKV("result.testEmailId", testEmailId)

	return s.createMailboxIfNotExists(ctx, span, tenant, mailboxAddress)
}

func (s *registrationService) createMailboxIfNotExists(ctx context.Context, span opentracing.Span, tenant, mailboxAddress string) error {
	mailbox, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, mailboxAddress)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get by mailbox"))
		return err
	}

	if mailbox == nil {
		mailboxRequest := AddMailboxRequest{
			Domain:          TEST_MAILBOX_DOMAIN,
			Username:        strings.ToLower(tenant),
			Password:        utils.GenerateLowerAlpha(1) + utils.GenerateKey(11, false),
			LinkedUserEmail: mailboxAddress,
			WebmailEnabled:  true,
			ForwardingTo:    []string{fmt.Sprintf("bcc@%s.customeros.ai", strings.ToLower(tenant))},
		}

		if err := s.services.MailboxService.AddMailbox(ctx, mailboxRequest); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to add mailbox"))
			return err
		}
	}

	return nil
}
