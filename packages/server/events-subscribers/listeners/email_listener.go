package listeners

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/security"
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

	emailValidationResponse, err := callApiValidateEmail(ctx, dependencies.CommonConfig.InternalServices.ValidationApiConfig, emailAddressToValidate)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error while calling email validation api"))
		return nil
	}

	if emailValidationResponse.Data.EmailData.SkippedValidation {
		span.LogFields(log.Bool("result.skippedValidation", true))
		dependencies.Logger.Warnf("Email %s for tenant %s skipped validation", emailId, tenant)
		return nil
	}
	if emailValidationResponse.Data.EmailData.RetryValidation {
		span.LogFields(log.Bool("result.retryValidation", true))
		dependencies.Logger.Warnf("Email %s for tenant %s need retry validation", emailId, tenant)
	}

	err = dependencies.CommonServices.EmailService.UpdateEmailValidationDetails(ctx, emailId, data_fields.EmailValidationFields{
		EmailAddress:      emailValidationResponse.Data.Syntax.CleanEmail,
		Domain:            emailValidationResponse.Data.Syntax.Domain,
		IsCatchAll:        emailValidationResponse.Data.DomainData.IsCatchAll,
		Deliverable:       emailValidationResponse.Data.EmailData.Deliverable,
		IsValidSyntax:     emailValidationResponse.Data.Syntax.IsValid,
		Username:          emailValidationResponse.Data.Syntax.User,
		ValidatedAt:       utils.Now(),
		IsRoleAccount:     emailValidationResponse.Data.EmailData.IsRoleAccount,
		IsSystemGenerated: emailValidationResponse.Data.EmailData.IsSystemGenerated,
		IsRisky: emailValidationResponse.Data.DomainData.IsFirewalled ||
			emailValidationResponse.Data.EmailData.IsRoleAccount ||
			emailValidationResponse.Data.EmailData.IsSystemGenerated ||
			emailValidationResponse.Data.EmailData.IsFreeAccount ||
			emailValidationResponse.Data.EmailData.IsMailboxFull ||
			!emailValidationResponse.Data.DomainData.IsPrimaryDomain,
		IsFirewalled:    emailValidationResponse.Data.DomainData.IsFirewalled,
		Provider:        emailValidationResponse.Data.DomainData.Provider,
		Firewall:        emailValidationResponse.Data.DomainData.SecureGatewayProvider,
		IsMailboxFull:   emailValidationResponse.Data.EmailData.IsMailboxFull,
		IsFreeAccount:   emailValidationResponse.Data.EmailData.IsFreeAccount,
		SmtpSuccess:     emailValidationResponse.Data.EmailData.SmtpSuccess,
		ResponseCode:    emailValidationResponse.Data.EmailData.ResponseCode,
		ErrorCode:       emailValidationResponse.Data.EmailData.ErrorCode,
		Description:     emailValidationResponse.Data.EmailData.Description,
		IsPrimaryDomain: emailValidationResponse.Data.DomainData.IsPrimaryDomain,
		PrimaryDomain:   emailValidationResponse.Data.DomainData.PrimaryDomain,
		AlternateEmail:  emailValidationResponse.Data.EmailData.AlternateEmail,
		RetryValidation: emailValidationResponse.Data.EmailData.RetryValidation,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to update email validation details"))
	}

	return err
}
