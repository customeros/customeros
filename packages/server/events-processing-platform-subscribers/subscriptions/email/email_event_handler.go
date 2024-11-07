package email

import (
	"bytes"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/email/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	validationmodel "github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"net/http"
)

type EmailEventHandler struct {
	services    *service.Services
	log         logger.Logger
	cfg         *config.Config
	caches      caches.Cache
	grpcClients *grpc_client.Clients
}

func NewEmailEventHandler(services *service.Services, log logger.Logger, cfg *config.Config, caches caches.Cache, grpcClients *grpc_client.Clients) *EmailEventHandler {
	emailEventHandler := EmailEventHandler{
		services:    services,
		log:         log,
		cfg:         cfg,
		caches:      caches,
		grpcClients: grpcClients,
	}

	return &emailEventHandler
}

func (h *EmailEventHandler) OnEmailValidate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailValidate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData event.EmailRequestValidationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	emailId := email.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	emailDbNode, err := h.services.CommonServices.Neo4jRepositories.EmailReadRepository.GetById(ctx, eventData.Tenant, emailId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	emailEntity := neo4jmapper.MapDbNodeToEmailEntity(emailDbNode)

	return h.validateEmail(ctx, eventData.Tenant, emailId, emailEntity.RawEmail)
}

func (h *EmailEventHandler) validateEmail(ctx context.Context, tenant, emailId, emailAddressToValidate string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.validateEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)

	emailValidationResponse, err := h.callApiValidateEmail(ctx, tenant, emailAddressToValidate)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error while calling email validation api"))
		return nil
	}

	if emailValidationResponse.Data.EmailData.SkippedValidation {
		span.LogFields(log.Bool("result.skippedValidation", true))
		h.log.Warnf("Email %s for tenant %s skipped validation", emailId, tenant)
		return nil
	}
	if emailValidationResponse.Data.EmailData.RetryValidation {
		span.LogFields(log.Bool("result.retryValidation", true))
		h.log.Warnf("Email %s for tenant %s need retry validation", emailId, tenant)
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
		request := emailpb.EmailValidationGrpcRequest{
			Tenant:        tenant,
			EmailId:       emailId,
			AppSource:     constants.AppSourceEventProcessingPlatformSubscribers,
			RawEmail:      emailAddressToValidate,
			Email:         emailValidationResponse.Data.Syntax.CleanEmail,
			Domain:        emailValidationResponse.Data.Syntax.Domain,
			Username:      emailValidationResponse.Data.Syntax.User,
			IsValidSyntax: emailValidationResponse.Data.Syntax.IsValid,
			IsRisky: emailValidationResponse.Data.DomainData.IsFirewalled ||
				emailValidationResponse.Data.EmailData.IsRoleAccount ||
				emailValidationResponse.Data.EmailData.IsSystemGenerated ||
				emailValidationResponse.Data.EmailData.IsFreeAccount ||
				emailValidationResponse.Data.EmailData.IsMailboxFull ||
				!emailValidationResponse.Data.DomainData.IsPrimaryDomain,
			IsFirewalled:      emailValidationResponse.Data.DomainData.IsFirewalled,
			Provider:          emailValidationResponse.Data.DomainData.Provider,
			Firewall:          emailValidationResponse.Data.DomainData.SecureGatewayProvider,
			IsCatchAll:        emailValidationResponse.Data.DomainData.IsCatchAll,
			Deliverable:       emailValidationResponse.Data.EmailData.Deliverable,
			IsMailboxFull:     emailValidationResponse.Data.EmailData.IsMailboxFull,
			IsRoleAccount:     emailValidationResponse.Data.EmailData.IsRoleAccount,
			IsSystemGenerated: emailValidationResponse.Data.EmailData.IsSystemGenerated,
			IsFreeAccount:     emailValidationResponse.Data.EmailData.IsFreeAccount,
			SmtpSuccess:       emailValidationResponse.Data.EmailData.SmtpSuccess,
			ResponseCode:      emailValidationResponse.Data.EmailData.ResponseCode,
			ErrorCode:         emailValidationResponse.Data.EmailData.ErrorCode,
			Description:       emailValidationResponse.Data.EmailData.Description,
			IsPrimaryDomain:   emailValidationResponse.Data.DomainData.IsPrimaryDomain,
			PrimaryDomain:     emailValidationResponse.Data.DomainData.PrimaryDomain,
			AlternateEmail:    emailValidationResponse.Data.EmailData.AlternateEmail,
			RetryValidation:   emailValidationResponse.Data.EmailData.RetryValidation,
		}
		return h.grpcClients.EmailClient.UpdateEmailValidation(ctx, &request)
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to call grpc"))
		h.log.Errorf("Failed sending email validation event for email %s for tenant %s: %s", emailId, tenant, err.Error())
	}
	return err
}

func (h *EmailEventHandler) callApiValidateEmail(ctx context.Context, tenant, emailAddress string) (*validationmodel.ValidateEmailResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.callApiValidateEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
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
	req, err := http.NewRequest("POST", h.cfg.Services.ValidationApi.Url+"/validateEmailV2", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request"))
		return nil, err
	}
	// Inject span context into the HTTP request
	req = commontracing.InjectSpanContextIntoHTTPRequest(req, span)

	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, h.cfg.Services.ValidationApi.ApiKey)
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
