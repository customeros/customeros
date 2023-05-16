package email_validation

import (
	"bytes"
	"context"
	"encoding/json"
	common_module "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
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

type EmailEventHandler struct {
	emailCommands *commands.EmailCommands
	log           logger.Logger
	cfg           *config.Config
}

type EmailValidate struct {
	Email string `json:"email" validate:"required,email"`
}

type EmailValidationResponseV1 struct {
	Error           string `json:"error"`
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

func (e *EmailEventHandler) OnEmailCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.EmailCreatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	emailValidate := EmailValidate{
		Email: strings.TrimSpace(eventData.RawEmail),
	}

	emailId := aggregate.GetEmailID(evt.AggregateID, eventData.Tenant)

	preValidationErr := validator.GetValidator().Struct(emailValidate)
	if preValidationErr != nil {
		e.emailCommands.FailEmailValidation.Handle(ctx, commands.NewFailedEmailValidationCommand(emailId, eventData.Tenant, preValidationErr.Error()))
	} else {
		evJSON, err := json.Marshal(emailValidate)
		if err != nil {
			tracing.TraceErr(span, err)
			e.emailCommands.FailEmailValidation.Handle(ctx, commands.NewFailedEmailValidationCommand(emailId, eventData.Tenant, err.Error()))
			return nil
		}
		requestBody := []byte(string(evJSON))
		req, err := http.NewRequest("POST", e.cfg.Services.ValidationApi+"/validateEmail", bytes.NewBuffer(requestBody))
		if err != nil {
			tracing.TraceErr(span, err)
			e.emailCommands.FailEmailValidation.Handle(ctx, commands.NewFailedEmailValidationCommand(emailId, eventData.Tenant, err.Error()))
			return nil
		}
		// Set the request headers
		req.Header.Set(common_module.ApiKeyHeader, e.cfg.Services.ValidationApiKey)
		req.Header.Set(common_module.TenantHeader, eventData.Tenant)

		// Make the HTTP request
		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			tracing.TraceErr(span, err)
			e.emailCommands.FailEmailValidation.Handle(ctx, commands.NewFailedEmailValidationCommand(emailId, eventData.Tenant, err.Error()))
			return nil
		}
		defer response.Body.Close()
		var result EmailValidationResponseV1
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			tracing.TraceErr(span, err)
			e.emailCommands.FailEmailValidation.Handle(ctx, commands.NewFailedEmailValidationCommand(emailId, eventData.Tenant, err.Error()))
			return nil
		}
		email := utils.StringFirstNonEmpty(result.Address, result.NormalizedEmail)
		e.emailCommands.EmailValidated.Handle(ctx, commands.NewEmailValidatedCommand(emailId, eventData.Tenant, emailValidate.Email,
			result.Error, result.Domain, result.Username, email, result.AcceptsMail, result.CanConnectSmtp,
			result.HasFullInbox, result.IsCatchAll, result.IsDisabled, result.IsValidSyntax))
	}

	return nil
}
