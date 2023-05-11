package email_validation

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
)

type PhoneNumberEventHandler struct {
	phoneNumberCommands *commands.PhoneNumberCommands
	log                 logger.Logger
	cfg                 *config.Config
}

type PhoneNumberValidate struct {
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}

type PhoneNumberValidationResponseV1 struct {
	Error       string `json:"error"`
	PhoneNumber string `json:"phoneNumber"`
}

func (h *PhoneNumberEventHandler) OnPhoneNumberCreate(ctx context.Context, evt eventstore.Event) error {
	//alexbalexb implement phone number validation
	//span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.OnEmailCreate")
	//defer span.Finish()
	//span.LogFields(log.String("AggregateID", evt.GetAggregateID()))
	//
	//var eventData events.EmailCreatedEvent
	//if err := evt.GetJsonData(&eventData); err != nil {
	//	tracing.TraceErr(span, err)
	//	return errors.Wrap(err, "evt.GetJsonData")
	//}
	//
	//emailValidate := EmailValidate{
	//	Email: strings.TrimSpace(eventData.RawEmail),
	//}
	//
	//emailId := aggregate.GetEmailID(evt.AggregateID, eventData.Tenant)
	//
	//preValidationErr := validator.GetValidator().Struct(emailValidate)
	//if preValidationErr != nil {
	//	e.emailCommands.FailEmailValidation.Handle(ctx, commands.NewFailEmailValidationCommand(emailId, eventData.Tenant, preValidationErr.Error()))
	//} else {
	//	evJSON, err := json.Marshal(emailValidate)
	//	if err != nil {
	//		tracing.TraceErr(span, err)
	//		e.emailCommands.FailEmailValidation.Handle(ctx, commands.NewFailEmailValidationCommand(emailId, eventData.Tenant, err.Error()))
	//		return nil
	//	}
	//	requestBody := []byte(string(evJSON))
	//	req, err := http.NewRequest("POST", e.cfg.Services.ValidationApi+"/validateEmail", bytes.NewBuffer(requestBody))
	//	if err != nil {
	//		tracing.TraceErr(span, err)
	//		e.emailCommands.FailEmailValidation.Handle(ctx, commands.NewFailEmailValidationCommand(emailId, eventData.Tenant, err.Error()))
	//		return nil
	//	}
	//	// Set the request headers
	//	req.Header.Set(common_module.ApiKeyHeader, e.cfg.Services.ValidationApiKey)
	//	req.Header.Set(common_module.TenantHeader, eventData.Tenant)
	//
	//	// Make the HTTP request
	//	client := &http.Client{}
	//	response, err := client.Do(req)
	//	if err != nil {
	//		tracing.TraceErr(span, err)
	//		e.emailCommands.FailEmailValidation.Handle(ctx, commands.NewFailEmailValidationCommand(emailId, eventData.Tenant, err.Error()))
	//		return nil
	//	}
	//	defer response.Body.Close()
	//	var result EmailValidationResponseV1
	//	err = json.NewDecoder(response.Body).Decode(&result)
	//	if err != nil {
	//		tracing.TraceErr(span, err)
	//		e.emailCommands.FailEmailValidation.Handle(ctx, commands.NewFailEmailValidationCommand(emailId, eventData.Tenant, err.Error()))
	//		return nil
	//	}
	//	e.emailCommands.EmailValidated.Handle(ctx, commands.NewEmailValidatedCommand(emailId, eventData.Tenant, emailValidate.Email,
	//		result.Error, result.Domain, result.Username, result.NormalizedEmail, result.AcceptsMail, result.CanConnectSmtp,
	//		result.HasFullInbox, result.IsCatchAll, result.IsDisabled, result.IsValidSyntax))
	//}

	return nil
}
