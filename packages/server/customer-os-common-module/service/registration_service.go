package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"strings"
)

type RegistrationService interface {
	PrepareDefaultTenantSetup(ctx context.Context, loggedInUserEmail string) error
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
	tenant := common.GetTenantFromContext(ctx)

	// Step 1 - Create default user
	testUserId, err := s.services.UserService.CreateTestUser(ctx, "Test", "Sender")
	span.LogKV("result.testUserId", testUserId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

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

	// Step 3 - Create postmark server for the tenant
	err = s.services.PostmarkService.CreateServer(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	// Step 4 - Register mailbox in opensrs
	password := utils.GenerateLowerAlpha(1) + utils.GenerateKey(11, false)
	username := strings.ToLower(tenant)
	forwardingTo := []string{fmt.Sprintf("bcc@%s.customeros.ai", strings.ToLower(tenant))}
	err = s.services.MailboxService.AddMailbox(ctx, TEST_MAILBOX_DOMAIN, username, password, mailboxAddress, true, true, forwardingTo)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	span.LogKV("result.mailbox", mailboxAddress)

	return nil
}
