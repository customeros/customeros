package phone_number_validation

import (
	"bytes"
	"context"
	"encoding/json"
	common_module "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type PhoneNumberEventHandler struct {
	repositories        *repository.Repositories
	phoneNumberCommands *commands.PhoneNumberCommands
	log                 logger.Logger
	cfg                 *config.Config
}

type PhoneNumberValidateRequest struct {
	PhoneNumber   string `json:"phoneNumber" validate:"required"`
	CountryCodeA2 string `json:"country"`
}

type PhoneNumberValidationResponseV1 struct {
	E164      string `json:"e164"`
	Error     string `json:"error"`
	Valid     bool   `json:"valid"`
	CountryA2 string `json:"countryA2"`
}

func (h *PhoneNumberEventHandler) OnPhoneNumberCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberEventHandler.OnPhoneNumberCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.PhoneNumberCreatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	tenant := eventData.Tenant
	rawPhoneNumber := eventData.RawPhoneNumber
	phoneNumberId := aggregate.GetPhoneNumberID(evt.AggregateID, tenant)

	countryCodeA2, err := h.repositories.PhoneNumberRepository.GetCountryCodeA2ForPhoneNumber(ctx, tenant, phoneNumberId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, err.Error())
		return nil
	}

	phoneNumberValidate := PhoneNumberValidateRequest{
		PhoneNumber:   strings.TrimSpace(eventData.RawPhoneNumber),
		CountryCodeA2: countryCodeA2,
	}

	preValidationErr := validator.GetValidator().Struct(phoneNumberValidate)
	if preValidationErr != nil {
		tracing.TraceErr(span, preValidationErr)
		h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, preValidationErr.Error())
		return nil
	} else {
		evJSON, err := json.Marshal(phoneNumberValidate)
		if err != nil {
			tracing.TraceErr(span, err)
			h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, preValidationErr.Error())
			return nil
		}
		requestBody := []byte(string(evJSON))
		req, err := http.NewRequest("POST", h.cfg.Services.ValidationApi+"/validatePhoneNumber", bytes.NewBuffer(requestBody))
		if err != nil {
			tracing.TraceErr(span, err)
			h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, preValidationErr.Error())
			return nil
		}
		// Set the request headers
		req.Header.Set(common_module.ApiKeyHeader, h.cfg.Services.ValidationApiKey)
		req.Header.Set(common_module.TenantHeader, tenant)

		// Make the HTTP request
		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			tracing.TraceErr(span, err)
			h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, preValidationErr.Error())
			return nil
		}
		defer response.Body.Close()
		var result PhoneNumberValidationResponseV1
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			tracing.TraceErr(span, err)
			h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, preValidationErr.Error())
			return nil
		}
		if !result.Valid {
			h.sendPhoneNumberFailedValidationEvent(ctx, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, result.Error)
		}
		h.phoneNumberCommands.PhoneNumberValidated.Handle(ctx, commands.NewPhoneNumberValidatedCommand(phoneNumberId, tenant, rawPhoneNumber, result.E164, result.CountryA2))
	}

	return nil
}

func (h *PhoneNumberEventHandler) sendPhoneNumberFailedValidationEvent(ctx context.Context, tenant, phoneNumberId, rawPhoneNumber, countryCodeA2, error string) {
	h.log.Errorf("Failed to validate phone number %s for tenant %s: %s", phoneNumberId, tenant, error)
	h.phoneNumberCommands.FailedPhoneNumberValidation.Handle(ctx, commands.NewFailedPhoneNumberValidationCommand(phoneNumberId, tenant, rawPhoneNumber, countryCodeA2, error))
}
