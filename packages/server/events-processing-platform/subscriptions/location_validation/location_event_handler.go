package location_validation

import (
	"bytes"
	"context"
	"encoding/json"
	common_module "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	utils_common "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/commands"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
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

type LocationEventHandler struct {
	repositories     *repository.Repositories
	locationCommands *commands.LocationCommands
	log              logger.Logger
	cfg              *config.Config
}

type LocationValidateRequest struct {
	Address       string `json:"address" validate:"required"`
	Country       string `json:"country"`
	International bool   `json:"international"`
}

type LocationValidationResponseV1 struct {
	Address *ValidatedAddress `json:"address"`
	Valid   bool              `json:"valid"`
	Error   *string           `json:"error"`
}
type ValidatedAddress struct {
	Country      string   `json:"country"`
	Region       string   `json:"region"`
	District     string   `json:"district"`
	Locality     string   `json:"locality"`
	Street       string   `json:"street"`
	Zip          string   `json:"zip"`
	PostalCode   string   `json:"postalCode"`
	AddressLine1 string   `json:"addressLine1"`
	AddressLine2 string   `json:"addressLine2"`
	AddressType  string   `json:"addressType"`
	HouseNumber  string   `json:"houseNumber"`
	PlusFour     string   `json:"plusFour"`
	Commercial   bool     `json:"commercial"`
	Predirection string   `json:"predirection"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	TimeZone     string   `json:"timeZone"`
	UtcOffset    int      `json:"utcOffset"`
}

func (h *LocationEventHandler) OnLocationCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationEventHandler.OnLocationCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.LocationCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed get event json: %v", err.Error())
		return errors.Wrap(err, "evt.GetJsonData")
	}

	tenant := eventData.Tenant
	locationId := aggregate.GetLocationObjectID(evt.AggregateID, tenant)

	if eventData.RawAddress == "" && eventData.LocationAddress.Address1 == "" && (eventData.LocationAddress.Street == "" || eventData.LocationAddress.HouseNumber == "") {
		return h.locationCommands.SkipLocationValidation.Handle(ctx, commands.NewSkippedLocationValidationCommand(locationId, tenant, "", "Missing raw Address"))
	}
	rawAddress := h.prepareRawAddress(eventData)
	country := h.prepareCountry(ctx, tenant, eventData.LocationAddress.Country)
	if country == "" {
		return h.locationCommands.SkipLocationValidation.Handle(ctx, commands.NewSkippedLocationValidationCommand(locationId, tenant, rawAddress, "Missing country"))
	}
	locationValidateRequest := LocationValidateRequest{
		Address:       rawAddress,
		Country:       country,
		International: !isCountryUSA(country),
	}

	preValidationErr := validator.GetValidator().Struct(locationValidateRequest)
	if preValidationErr != nil {
		tracing.TraceErr(span, preValidationErr)
		h.log.Errorf("Failed to pre-validate location: %v", preValidationErr.Error())
		return h.sendLocationFailedValidationEvent(ctx, tenant, locationId, rawAddress, country, preValidationErr.Error())
	}
	evJSON, err := json.Marshal(locationValidateRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to marshal location validation request: %v", err.Error())
		return h.sendLocationFailedValidationEvent(ctx, tenant, locationId, rawAddress, country, err.Error())
	}
	requestBody := []byte(string(evJSON))
	req, err := http.NewRequest("POST", h.cfg.Services.ValidationApi+"/validateAddress", bytes.NewBuffer(requestBody))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to create location validation request: %v", err.Error())
		return h.sendLocationFailedValidationEvent(ctx, tenant, locationId, rawAddress, country, err.Error())
	}
	// Set the request headers
	req.Header.Set(common_module.ApiKeyHeader, h.cfg.Services.ValidationApiKey)
	req.Header.Set(common_module.TenantHeader, tenant)

	// Make the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to send location validation request: %v", err.Error())
		return h.sendLocationFailedValidationEvent(ctx, tenant, locationId, rawAddress, country, err.Error())
	}
	defer response.Body.Close()
	var result LocationValidationResponseV1
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to decode location validation response: %v", err.Error())
		return h.sendLocationFailedValidationEvent(ctx, tenant, locationId, rawAddress, country, err.Error())
	}
	if !result.Valid {
		return h.sendLocationFailedValidationEvent(ctx, tenant, locationId, rawAddress, country, *result.Error)
	}

	locationAddressFields := models.LocationAddressFields{
		Country:      result.Address.Country,
		Region:       result.Address.Region,
		District:     result.Address.District,
		Locality:     result.Address.Locality,
		Street:       result.Address.Street,
		Address1:     result.Address.AddressLine1,
		Address2:     result.Address.AddressLine2,
		Zip:          result.Address.Zip,
		AddressType:  result.Address.AddressType,
		HouseNumber:  result.Address.HouseNumber,
		PostalCode:   result.Address.PostalCode,
		PlusFour:     result.Address.PlusFour,
		Commercial:   result.Address.Commercial,
		Predirection: result.Address.Predirection,
		Latitude:     result.Address.Latitude,
		Longitude:    result.Address.Longitude,
		TimeZone:     result.Address.TimeZone,
		UtcOffset:    result.Address.UtcOffset,
	}
	return h.locationCommands.LocationValidated.Handle(ctx, commands.NewLocationValidatedCommand(locationId, tenant, rawAddress, country, locationAddressFields))
}

func (h *LocationEventHandler) prepareRawAddress(eventData events.LocationCreateEvent) string {
	rawAddress := strings.TrimSpace(eventData.RawAddress)
	if rawAddress == "" {
		rawAddress = constructRawAddressForValidationFromLocationAddressFields(eventData)
	}
	return rawAddress
}

func (h *LocationEventHandler) prepareCountry(ctx context.Context, tenant, eventCountry string) string {
	if eventCountry != "" {
		return eventCountry
	}
	country, err := h.repositories.CountryRepository.GetDefaultCountryCodeA3(ctx, tenant)
	if err != nil {
		return ""
	}
	return country
}

func constructRawAddressForValidationFromLocationAddressFields(eventData events.LocationCreateEvent) string {
	rawAddress :=
		eventData.LocationAddress.HouseNumber + " " +
			eventData.LocationAddress.Street + " " +
			eventData.LocationAddress.Address1 + " " +
			eventData.LocationAddress.Address2 + " " +
			utils_common.StringFirstNonEmpty(eventData.LocationAddress.Zip, eventData.LocationAddress.PostalCode) + ", " +
			eventData.LocationAddress.Locality
	if eventData.LocationAddress.Locality != "" {
		rawAddress += ","
	}
	rawAddress += " " + eventData.LocationAddress.District + " " +
		eventData.LocationAddress.Region
	return rawAddress
}

func isCountryUSA(country string) bool {
	return country == "USA" || country == "US" || country == "United States" || country == "United States of America"
}

func (h *LocationEventHandler) sendLocationFailedValidationEvent(ctx context.Context, tenant, locationId, rawAddress, country, error string) error {
	h.log.Errorf("Failed validating location %s for tenant %s: %s", locationId, tenant, error)
	return h.locationCommands.FailedLocationValidation.Handle(ctx, commands.NewFailedLocationValidationCommand(locationId, tenant, rawAddress, country, error))
}
