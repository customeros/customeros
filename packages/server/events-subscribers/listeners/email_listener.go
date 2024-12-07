package listeners

import (
	"bytes"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"net/http"
)

func OnRequestedValidateEmail(ctx context.Context, services *service.Services, input any) error {
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

	emailDbNode, err := services.Neo4jRepositories.EmailReadRepository.GetById(ctx, tenant, emailId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Failed to get email node from neo4j"))
		return err
	}
	emailEntity := neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)

	return validateEmail(ctx, services, emailId, emailEntity.RawEmail)
}

func validateEmail(ctx context.Context, services *service.Services, emailId, emailAddressToValidate string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.validateEmail")
	defer span.Finish()
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, emailId)

	emailValidationResponse, err := callApiValidateEmail(ctx, services.GlobalConfig.InternalServices.ValidationApiConfig, emailAddressToValidate)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error while calling email validation api"))
		return nil
	}

	if emailValidationResponse.Data.EmailData.SkippedValidation {
		span.LogFields(log.Bool("result.skippedValidation", true))
		services.Logger.Warnf("Email %s for tenant %s skipped validation", emailId, tenant)
		return nil
	}
	if emailValidationResponse.Data.EmailData.RetryValidation {
		span.LogFields(log.Bool("result.retryValidation", true))
		services.Logger.Warnf("Email %s for tenant %s need retry validation", emailId, tenant)
	}

	err = services.EmailService.UpdateEmailValidationDetails(ctx, emailId, data_fields.EmailValidationFields{
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

func callApiValidateEmail(ctx context.Context, validationApiConfig config.ValidationAPIConfig, emailAddress string) (*validationmodel.ValidateEmailResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.callApiValidateEmail")
	defer span.Finish()
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)
	span.LogKV("emailAddress", emailAddress)

	// prepare validation api request
	requestJSON, err := json.Marshal(validationmodel.ValidateEmailRequestWithOptions{
		Email: emailAddress,
		Options: validationmodel.ValidateEmailRequestOptions{
			VerifyCatchAll: true,
		},
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal request"))
		return nil, err
	}
	requestBody := []byte(string(requestJSON))
	req, err := http.NewRequest("POST", validationApiConfig.Url+"/validateEmailV2", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = tracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, validationApiConfig.ApiKey)
	req.Header.Set(security.TenantHeader, tenant)

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to perform request"))
		return nil, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to read response body"))
		return nil, err
	}

	var validationResponse validationmodel.ValidateEmailResponse
	err = json.Unmarshal(responseBody, &validationResponse)
	if err != nil {
		span.LogFields(log.String("response.body", string(responseBody)))
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode response"))
		return nil, err
	}
	if validationResponse.Data == nil {
		err = errors.New("email validation response data is empty: " + validationResponse.InternalMessage)
		tracing.TraceErr(span, err)
		return nil, err
	}
	tracing.LogObjectAsJson(span, "email validation response", validationResponse)
	return &validationResponse, nil
}
