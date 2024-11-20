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

func (s *registrationService) PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegistrationService.PrepareDefaultTenantSetup")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("loggedInUserEmail", loggedInUserEmail)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.ConfigureTestMailbox(ctx)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error configuring test mailbox during tenant onboarding"))
	}

	err = s.CreatePostmarkServer(ctx)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating postmark server during tenant onboarding"))
	}

	return nil
}

func (s *registrationService) ConfigureTestMailbox(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegistrationService.ConfigureTestMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	var testUserId string

	existingTestUser, err := s.services.Neo4jRepositories.UserReadRepository.FindTestUser(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if existingTestUser == nil {
		testUserId, err = s.services.UserService.CreateUser(ctx, neo4jentity.UserEntity{
			FirstName: "Test",
			LastName:  "Sender",
			Test:      true,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		testUserId = mapper.MapDbNodeToUserEntity(existingTestUser).Id
	}
	span.LogKV("result.testUserId", testUserId)

	// Step 2 - Create test email node for the user
	mailboxAddress := strings.ToLower(fmt.Sprintf("%s@%s", tenant, TEST_MAILBOX_DOMAIN))
	testEmailId, err := s.services.EmailService.Merge(ctx, tenant, EmailFields{
		Email: mailboxAddress,
	}, &LinkWith{
		Type: model.USER,
		Id:   testUserId,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	span.LogKV("result.testEmailId", testEmailId)

	mailbox, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, mailboxAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if mailbox == nil {
		// Step 4 - Register mailbox in opensrs
		var mailboxRequest MailboxRequest

		mailboxRequest.Domain = TEST_MAILBOX_DOMAIN
		mailboxRequest.Username = strings.ToLower(tenant)
		mailboxRequest.Password = utils.GenerateLowerAlpha(1) + utils.GenerateKey(11, false)
		mailboxRequest.LinkedUserEmail = mailboxAddress
		mailboxRequest.WebmailEnabled = true
		mailboxRequest.ForwardingEnabled = true
		mailboxRequest.ForwardingTo = []string{fmt.Sprintf("bcc@%s.customeros.ai", strings.ToLower(tenant))}

		err = s.services.MailboxService.AddMailbox(ctx, mailboxRequest)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	span.LogKV("result.mailboxAddress", mailboxAddress)

	return nil
}

func (s *registrationService) CreatePostmarkServer(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegistrationService.CreatePostmarkServer")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Step 3 - Create postmark server for the tenant
	err = s.services.PostmarkService.CreateServer(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return nil
}
