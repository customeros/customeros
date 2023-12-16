package email_validation

import (
	"bytes"
	"context"
	"encoding/json"
	common_module "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type emailEventHandler struct {
	emailCommands *command_handler.CommandHandlers
	log           logger.Logger
	cfg           *config.Config
	grpcClients   *grpc_client.Clients
}

func NewEmailEventHandler(emailCommands *command_handler.CommandHandlers, log logger.Logger, cfg *config.Config, grpcClients *grpc_client.Clients) *emailEventHandler {
	return &emailEventHandler{
		emailCommands: emailCommands,
		log:           log,
		cfg:           cfg,
		grpcClients:   grpcClients,
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

	emailValidate := EmailValidate{
		Email: strings.TrimSpace(eventData.RawEmail),
	}

	emailId := aggregate.GetEmailObjectID(evt.AggregateID, eventData.Tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	preValidationErr := validator.GetValidator().Struct(emailValidate)
	if preValidationErr != nil {
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, preValidationErr.Error(), span)
	}
	evJSON, err := json.Marshal(emailValidate)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, err.Error(), span)
	}
	requestBody := []byte(string(evJSON))
	req, err := http.NewRequest("POST", h.cfg.Services.ValidationApi+"/validateEmail", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, err.Error(), span)
	}
	// Set the request headers
	req.Header.Set(common_module.ApiKeyHeader, h.cfg.Services.ValidationApiKey)
	req.Header.Set(common_module.TenantHeader, eventData.Tenant)

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, err.Error(), span)
	}
	defer response.Body.Close()
	var result EmailValidationResponseV1
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, err.Error(), span)
	}
	if result.IsReachable == "" {
		errMsg := utils.StringFirstNonEmpty(result.Error, "IsReachable flag not set. Email not passed validation.")
		return h.sendEmailFailedValidationEvent(ctx, eventData.Tenant, emailId, errMsg, span)
	}
	email := utils.StringFirstNonEmpty(result.Address, result.NormalizedEmail)
	return h.emailCommands.EmailValidated.Handle(ctx, command.NewEmailValidatedCommand(emailId, eventData.Tenant, emailValidate.Email, result.IsReachable,
		result.Error, result.Domain, result.Username, email, result.AcceptsMail, result.CanConnectSmtp,
		result.HasFullInbox, result.IsCatchAll, result.IsDisabled, result.IsValidSyntax))
}

func (h *emailEventHandler) sendEmailFailedValidationEvent(ctx context.Context, tenant, emailId string, errorMessage string, span opentracing.Span) error {
	h.log.Errorf("Failed validating email %s for tenant %s: %s", emailId, tenant, errorMessage)
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := h.grpcClients.EmailClient.FailEmailValidation(ctx, &emailpb.FailEmailValidationGrpcRequest{
		Tenant:       tenant,
		EmailId:      emailId,
		AppSource:    constants.AppSourceEventProcessingPlatform,
		ErrorMessage: utils.StringFirstNonEmpty(errorMessage, "Error message not available"),
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed sending failed email validation event for email %s for tenant %s: %s", emailId, tenant, err.Error())
	}
	return err
}
