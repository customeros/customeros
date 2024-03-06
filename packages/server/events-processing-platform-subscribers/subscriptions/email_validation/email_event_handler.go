package email_validation

import (
	"bytes"
	"context"
	"encoding/json"
	common_module "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type emailEventHandler struct {
	log         logger.Logger
	cfg         *config.Config
	grpcClients *grpc_client.Clients
}

func NewEmailEventHandler(log logger.Logger, cfg *config.Config, grpcClients *grpc_client.Clients) *emailEventHandler {
	return &emailEventHandler{
		log:         log,
		cfg:         cfg,
		grpcClients: grpcClients,
	}
}

type EmailValidate struct {
	Email string `json:"email" validate:"required,email"`
}

type EmailValidationResponseV1 struct {
	Error           string `json:"error"`
	IsReachable     string `json:"isReachable"`
	Email           string `json:"email"`
	AcceptsMail     bool   `json:"acceptsMail"`
	CanConnectSmtp  bool   `json:"canConnectSmtp"`
	HasFullInbox    bool   `json:"hasFullInbox"`
	IsCatchAll      bool   `json:"isCatchAll"`
	IsDeliverable   bool   `json:"isDeliverable"`
	IsDisabled      bool   `json:"isDisabled"`
	Address         string `json:"address"`
	Domain          string `json:"domain"`
	IsValidSyntax   bool   `json:"isValidSyntax"`
	Username        string `json:"username"`
	NormalizedEmail string `json:"normalizedEmail"`
}

func (h *emailEventHandler) ValidateEmail(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.ValidateEmail")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.EmailCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}
	span.SetTag(tracing.SpanTagTenant, eventData.Tenant)

	emailValidate := EmailValidate{
		Email: strings.TrimSpace(eventData.RawEmail),
	}

	emailId := aggregate.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	preValidationErr := validator.GetValidator().Struct(emailValidate)
	if preValidationErr != nil {
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, preValidationErr.Error())
	}
	evJSON, err := json.Marshal(emailValidate)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, err.Error())
	}
	requestBody := []byte(string(evJSON))
	req, err := http.NewRequest("POST", h.cfg.Services.ValidationApi+"/validateEmail", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, err.Error())
	}
	// Set the request headers
	req.Header.Set(common_module.ApiKeyHeader, h.cfg.Services.ValidationApiKey)
	req.Header.Set(common_module.TenantHeader, eventData.Tenant)

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, err.Error())
	}
	defer response.Body.Close()
	var result EmailValidationResponseV1
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, err.Error())
	}
	if result.IsReachable == "" {
		errMsg := utils.StringFirstNonEmpty(result.Error, "IsReachable flag not set. Email not passed validation.")
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, errMsg)
	}
	email := utils.StringFirstNonEmpty(result.Address, result.NormalizedEmail)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
		return h.grpcClients.EmailClient.PassEmailValidation(ctx, &emailpb.PassEmailValidationGrpcRequest{
			Tenant:         eventData.Tenant,
			EmailId:        emailId,
			AppSource:      constants.AppSourceEventProcessingPlatform,
			RawEmail:       emailValidate.Email,
			Email:          email,
			IsReachable:    result.IsReachable,
			ErrorMessage:   result.Error,
			Domain:         result.Domain,
			Username:       result.Username,
			AcceptsMail:    result.AcceptsMail,
			CanConnectSmtp: result.CanConnectSmtp,
			HasFullInbox:   result.HasFullInbox,
			IsCatchAll:     result.IsCatchAll,
			IsDisabled:     result.IsDisabled,
			IsValidSyntax:  result.IsValidSyntax,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed sending passed email validation event for email %s for tenant %s: %s", emailId, eventData.Tenant, err.Error())
	}
	return err
}

func (h *emailEventHandler) sendEmailFailedValidationEvent(ctx context.Context, tenant, emailId string, errorMessage string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.sendEmailFailedValidationEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("emailId", emailId), log.String("errorMessage", errorMessage))

	h.log.Errorf("Failed validating email %s for tenant %s: %s", emailId, tenant, errorMessage)
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
		return h.grpcClients.EmailClient.FailEmailValidation(ctx, &emailpb.FailEmailValidationGrpcRequest{
			Tenant:       tenant,
			EmailId:      emailId,
			AppSource:    constants.AppSourceEventProcessingPlatform,
			ErrorMessage: utils.StringFirstNonEmpty(errorMessage, "Error message not available"),
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed sending failed email validation event for email %s for tenant %s: %s", emailId, tenant, err.Error())
	}
	return err
}
