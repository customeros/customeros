package location

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	aiConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/config"
	ai "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commonTracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/constants"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	locationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/location"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

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

type LocationEventHandler struct {
	repositories *repository.Repositories
	log          logger.Logger
	cfg          *config.Config
	grpcClients  *grpc_client.Clients
	aiModel      ai.AiModel
}

func NewLocationEventHandler(repositories *repository.Repositories, log logger.Logger, cfg *config.Config, grpcClients *grpc_client.Clients) *LocationEventHandler {
	return &LocationEventHandler{
		repositories: repositories,
		log:          log,
		cfg:          cfg,
		grpcClients:  grpcClients,
		aiModel: ai.NewAiModel(ai.AnthropicModelType, aiConfig.Config{
			Anthropic: aiConfig.AiModelConfigAnthropic{
				ApiPath: cfg.Services.Ai.ApiPath,
				ApiKey:  cfg.Services.Ai.ApiKey,
				Model:   constants.AnthropicApiModel,
			},
		}),
	}
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
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*locationpb.LocationIdGrpcResponse](func() (*locationpb.LocationIdGrpcResponse, error) {
			return h.grpcClients.LocationClient.SkipLocationValidation(ctx, &locationpb.SkipLocationValidationGrpcRequest{
				Tenant:     tenant,
				LocationId: locationId,
				AppSource:  constants.AppSourceEventProcessingPlatformSubscribers,
				RawAddress: "",
				Reason:     "Missing raw Address",
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed sending skipped location validation event for location %s for tenant %s: %s", locationId, tenant, err.Error())
		}
		return err
	}
	rawAddress := h.prepareRawAddress(eventData)
	country := h.prepareCountry(ctx, tenant, eventData.LocationAddress.Country)
	if country == "" {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*locationpb.LocationIdGrpcResponse](func() (*locationpb.LocationIdGrpcResponse, error) {
			return h.grpcClients.LocationClient.SkipLocationValidation(ctx, &locationpb.SkipLocationValidationGrpcRequest{
				Tenant:     tenant,
				LocationId: locationId,
				AppSource:  constants.AppSourceEventProcessingPlatformSubscribers,
				RawAddress: rawAddress,
				Reason:     "Missing country",
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed sending skipped location validation event for location %s for tenant %s: %s", locationId, tenant, err.Error())
		}
		return err
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
	// Inject span context into the HTTP request
	req = commonTracing.InjectSpanContextIntoHTTPRequest(req, span)
	// Set the request headers
	req.Header.Set(security.ApiKeyHeader, h.cfg.Services.ValidationApiKey)
	req.Header.Set(security.TenantHeader, tenant)

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
		return h.sendLocationFailedValidationEvent(ctx, tenant, locationId, rawAddress, country, utils.IfNotNilStringWithDefault(result.Error, "missing error in validation response"))
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = subscriptions.CallEventsPlatformGRPCWithRetry[*locationpb.LocationIdGrpcResponse](func() (*locationpb.LocationIdGrpcResponse, error) {
		request := locationpb.PassLocationValidationGrpcRequest{
			Tenant:       tenant,
			LocationId:   locationId,
			AppSource:    constants.AppSourceEventProcessingPlatformSubscribers,
			RawAddress:   rawAddress,
			Country:      country,
			Region:       result.Address.Region,
			District:     result.Address.District,
			Locality:     result.Address.Locality,
			Street:       result.Address.Street,
			AddressLine1: result.Address.AddressLine1,
			AddressLine2: result.Address.AddressLine2,
			ZipCode:      result.Address.Zip,
			PostalCode:   result.Address.PostalCode,
			AddressType:  result.Address.AddressType,
			HouseNumber:  result.Address.HouseNumber,
			PlusFour:     result.Address.PlusFour,
			Commercial:   result.Address.Commercial,
			Predirection: result.Address.Predirection,
			TimeZone:     result.Address.TimeZone,
			UtcOffset:    int32(result.Address.UtcOffset),
		}
		if result.Address.Latitude != nil {
			request.Latitude = utils.FloatToString(result.Address.Latitude)
		}
		if result.Address.Longitude != nil {
			request.Longitude = utils.FloatToString(result.Address.Longitude)
		}
		return h.grpcClients.LocationClient.PassLocationValidation(ctx, &request)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed sending skipped location validation event for location %s for tenant %s: %s", locationId, tenant, err.Error())
	}
	return err
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
	country, err := h.repositories.Neo4jRepositories.CountryReadRepository.GetDefaultCountryCodeA3(ctx, tenant)
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
			utils.StringFirstNonEmpty(eventData.LocationAddress.Zip, eventData.LocationAddress.PostalCode) + ", " +
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

func (h *LocationEventHandler) sendLocationFailedValidationEvent(ctx context.Context, tenant, locationId, rawAddress, country, errorMessage string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailEventHandler.sendEmailFailedValidationEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("locationId", locationId), log.String("errorMessage", errorMessage))

	h.log.Errorf("Failed validating location %s for tenant %s: %s", locationId, tenant, errorMessage)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := subscriptions.CallEventsPlatformGRPCWithRetry[*locationpb.LocationIdGrpcResponse](func() (*locationpb.LocationIdGrpcResponse, error) {
		return h.grpcClients.LocationClient.FailLocationValidation(ctx, &locationpb.FailLocationValidationGrpcRequest{
			Tenant:       tenant,
			LocationId:   locationId,
			AppSource:    constants.AppSourceEventProcessingPlatformSubscribers,
			RawAddress:   rawAddress,
			Country:      country,
			ErrorMessage: utils.StringFirstNonEmpty(errorMessage, "Error message not available"),
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed sending failed location validation event for location %s for tenant %s: %s", locationId, tenant, err.Error())
	}
	return err
}

func (h *LocationEventHandler) ExtractAndEnrichLocation(ctx context.Context, tenant, address string) (*Location, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationEventHandler.ExtractAndEnrichLocation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("address", address))

	if strings.TrimSpace(address) == "" {
		return nil, errors.New("address is empty")
	}

	// Step 1: Check if mapping exists
	locationMapping, err := h.repositories.PostgresRepositories.AiLocationMappingRepository.GetLatestLocationMappingByInput(ctx, address)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get location mapping"))
	}
	if locationMapping != nil {
		var location Location
		err = json.Unmarshal([]byte(locationMapping.ResponseJson), &location)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal location"))
			return nil, err
		}
		return &location, nil
	}

	// Step 2: Use AI to enrich the location
	prompt := fmt.Sprintf(h.cfg.Services.Anthropic.LocationEnrichmentPrompt, address)
	promptLog := postgresEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      constants.AppSourceEventProcessingPlatformSubscribers,
		Provider:       constants.Anthropic,
		Model:          constants.AnthropicApiModel,
		PromptType:     constants.PromptTypeExtractLocationValue,
		Tenant:         &tenant,
		PromptTemplate: &h.cfg.Services.Anthropic.LocationEnrichmentPrompt,
		Prompt:         prompt,
	}
	promptStoreLogId, err := h.repositories.PostgresRepositories.AiPromptLogRepository.Store(promptLog)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error storing prompt log: %v", err)
	}

	aiResult, err := h.aiModel.Inference(ctx, prompt)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get AI response"))
		h.log.Errorf("Error invoking AI: %s", err.Error())
		storeErr := h.repositories.PostgresRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, errors.Wrap(storeErr, "failed to update prompt log with error"))
			h.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
		return nil, err
	} else {
		storeErr := h.repositories.PostgresRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId, aiResult)
		if storeErr != nil {
			tracing.TraceErr(span, errors.Wrap(storeErr, "failed to update prompt log with ai response"))
			h.log.Errorf("Error updating prompt log with ai response: %v", storeErr)
		}
	}

	var location Location
	err = json.Unmarshal([]byte(aiResult), &location)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal location"))
		return nil, err
	}

	// Step 3: Store the mapping
	locationMapping = &postgresEntity.AiLocationMapping{
		Input:         address,
		ResponseJson:  aiResult,
		AiPromptLogId: promptStoreLogId,
	}
	err = h.repositories.PostgresRepositories.AiLocationMappingRepository.AddLocationMapping(ctx, *locationMapping)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to store location mapping"))
		h.log.Errorf("Error storing location mapping: %v", err)
	}

	return &location, nil
}
