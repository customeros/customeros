package email_validation

import (
	"bytes"
	"context"
	"encoding/json"
	common_module "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
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
	emailCommands *command_handler.EmailCommands
	log           logger.Logger
	cfg           *config.Config
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

	preValidationErr := validator.GetValidator().Struct(emailValidate)
	if preValidationErr != nil {
		return h.emailCommands.FailEmailValidation.Handle(ctx, command.NewFailedEmailValidationCommand(emailId, eventData.Tenant, preValidationErr.Error()))
	}
	evJSON, err := json.Marshal(emailValidate)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, emailId, eventData.Tenant, err.Error())
	}
	requestBody := []byte(string(evJSON))
	req, err := http.NewRequest("POST", h.cfg.Services.ValidationApi+"/validateEmail", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, emailId, eventData.Tenant, err.Error())
	}
	// Set the request headers
	req.Header.Set(common_module.ApiKeyHeader, h.cfg.Services.ValidationApiKey)
	req.Header.Set(common_module.TenantHeader, eventData.Tenant)

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, emailId, eventData.Tenant, err.Error())
	}
	defer response.Body.Close()
	var result EmailValidationResponseV1
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		tracing.TraceErr(span, err)
		return h.sendEmailFailedValidationEvent(ctx, emailId, eventData.Tenant, err.Error())
	}
	if result.IsReachable == "" {
		errMsg := utils.StringFirstNonEmpty(result.Error, "IsReachable flag not set. Email not passed validation.")
		return h.sendEmailFailedValidationEvent(ctx, emailId, eventData.Tenant, errMsg)
	}
	email := utils.StringFirstNonEmpty(result.Address, result.NormalizedEmail)
	return h.emailCommands.EmailValidated.Handle(ctx, command.NewEmailValidatedCommand(emailId, eventData.Tenant, emailValidate.Email, result.IsReachable,
		result.Error, result.Domain, result.Username, email, result.AcceptsMail, result.CanConnectSmtp,
		result.HasFullInbox, result.IsCatchAll, result.IsDisabled, result.IsValidSyntax))
}

func (h *emailEventHandler) sendEmailFailedValidationEvent(ctx context.Context, emailId, tenant string, errMsg string) error {
	h.log.Errorf("Failed validating email %s for tenant %s: %s", emailId, tenant, errMsg)
	return h.emailCommands.FailEmailValidation.Handle(ctx, command.NewFailedEmailValidationCommand(emailId, tenant, errMsg))
}
