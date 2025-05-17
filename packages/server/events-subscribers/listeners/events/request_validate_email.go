package events_listeners

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type RequestValidateEmailListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewRequestValidateEmailListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &RequestValidateEmailListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.RequestValidateEmail](), // subscribed event
			events.QueueEvents, // listening on CustomerOS Events queue
		),
		dependencies: deps,
	}
}

func (l *RequestValidateEmailListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartListenerSpan(ctx, "RequestValidateEmailListener.Handle")
	defer spans.Finish()
	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	emailId := event.Event.EntityId
	spans.LogKV("emailId", emailId)

	emailDbNode, err := l.dependencies.Neo4jRepositories.EmailReadRepository.GetById(ctx, event.Event.Tenant, emailId)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "Failed to get email node from neo4j"))
		return err
	}
	emailEntity := neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)

	return l.validateEmail(ctx, emailId, emailEntity.RawEmail)
}

func (l *RequestValidateEmailListener) validateEmail(ctx context.Context, emailId, emailAddressToValidate string) error {
	spans, ctx := telemetry.StartListenerSpan(ctx, "RequestValidateEmailListener.validateEmail")
	defer spans.Finish()
	tenant := common.GetTenantFromContext(ctx)
	spans.TagTenant(tenant)
	spans.TagEntity(emailId)

	emailValidationResponse, err := l.dependencies.CommonServices.VerifyService.ValidateEmail(
		ctx, emailAddressToValidate)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "error while calling email validation api"))
		return nil
	}

	if emailValidationResponse.EmailData.SkippedValidation {
		spans.LogKV("result.skippedValidation", true)
		l.dependencies.Logger.Warnf("Email %s for tenant %s skipped validation", emailId, tenant)
		return nil
	}
	if emailValidationResponse.EmailData.RetryValidation {
		spans.LogKV("result.retryValidation", true)
		l.dependencies.Logger.Warnf("Email %s for tenant %s need retry validation", emailId, tenant)
	}

	err = l.dependencies.CommonServices.EmailService.UpdateEmailValidationDetails(ctx, emailId, data_fields.EmailValidationFields{
		EmailAddress:      emailValidationResponse.Syntax.CleanEmail,
		Domain:            emailValidationResponse.Syntax.Domain,
		IsCatchAll:        emailValidationResponse.DomainData.IsCatchAll,
		Deliverable:       emailValidationResponse.EmailData.Deliverable,
		IsValidSyntax:     emailValidationResponse.Syntax.IsValid,
		Username:          emailValidationResponse.Syntax.User,
		ValidatedAt:       utils.Now(),
		IsRoleAccount:     emailValidationResponse.EmailData.IsRoleAccount,
		IsSystemGenerated: emailValidationResponse.EmailData.IsSystemGenerated,
		IsRisky: emailValidationResponse.DomainData.IsFirewalled ||
			emailValidationResponse.EmailData.IsRoleAccount ||
			emailValidationResponse.EmailData.IsSystemGenerated ||
			emailValidationResponse.EmailData.IsFreeAccount ||
			emailValidationResponse.EmailData.IsMailboxFull ||
			!emailValidationResponse.DomainData.IsPrimaryDomain,
		IsFirewalled:    emailValidationResponse.DomainData.IsFirewalled,
		Provider:        emailValidationResponse.DomainData.Provider,
		Firewall:        emailValidationResponse.DomainData.SecureGatewayProvider,
		IsMailboxFull:   emailValidationResponse.EmailData.IsMailboxFull,
		IsFreeAccount:   emailValidationResponse.EmailData.IsFreeAccount,
		SmtpSuccess:     emailValidationResponse.EmailData.SmtpSuccess,
		ResponseCode:    emailValidationResponse.EmailData.ResponseCode,
		ErrorCode:       emailValidationResponse.EmailData.ErrorCode,
		Description:     emailValidationResponse.EmailData.Description,
		IsPrimaryDomain: emailValidationResponse.DomainData.IsPrimaryDomain,
		PrimaryDomain:   emailValidationResponse.DomainData.PrimaryDomain,
		AlternateEmail:  emailValidationResponse.EmailData.AlternateEmail,
		RetryValidation: emailValidationResponse.EmailData.RetryValidation,
	})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to update email validation details"))
	}

	return err
}
