package listeners

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func OnRequestedValidateEmail(ctx context.Context, dependencies *model.DependencyContainer, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnRequestedValidateEmail")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message := input.(*dto.Event)
	emailId := message.Event.EntityId
	span.SetTag(tracing.SpanTagEntityId, emailId)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Missing tenant in email request validation event")
		tracing.TraceErr(span, err)
		return err
	}

	emailDbNode, err := dependencies.Neo4jRepositories.EmailReadRepository.GetById(ctx, tenant, emailId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Failed to get email node from neo4j"))
		return err
	}
	emailEntity := neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)

	return validateEmail(ctx, dependencies, emailId, emailEntity.RawEmail)
}

func validateEmail(ctx context.Context, dependencies *model.DependencyContainer, emailId, emailAddressToValidate string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.validateEmail")
	defer span.Finish()
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, emailId)

	emailValidationResponse, err := dependencies.CommonServices.VerifyService.ValidateEmail(
		ctx, emailAddressToValidate)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error while calling email validation api"))
		return nil
	}

	if emailValidationResponse.EmailData.SkippedValidation {
		span.LogFields(log.Bool("result.skippedValidation", true))
		dependencies.Logger.Warnf("Email %s for tenant %s skipped validation", emailId, tenant)
		return nil
	}
	if emailValidationResponse.EmailData.RetryValidation {
		span.LogFields(log.Bool("result.retryValidation", true))
		dependencies.Logger.Warnf("Email %s for tenant %s need retry validation", emailId, tenant)
	}

	err = dependencies.CommonServices.EmailService.UpdateEmailValidationDetails(ctx, emailId, data_fields.EmailValidationFields{
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
		tracing.TraceErr(span, errors.Wrap(err, "failed to update email validation details"))
	}

	return err
}
