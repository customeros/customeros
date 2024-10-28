package service

import (
	"encoding/json"
	"fmt"
	aiConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/config"
	ai "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"strings"
)

type Location struct {
	Country       string   `json:"country"`
	CountryCodeA2 string   `json:"countryCodeA2"`
	CountryCodeA3 string   `json:"countryCodeA3"`
	Region        string   `json:"region"`
	Locality      string   `json:"locality"`
	Address       string   `json:"address"`
	Address2      string   `json:"address2"`
	Zip           string   `json:"zip"`
	AddressType   string   `json:"addressType"`
	HouseNumber   string   `json:"houseNumber"`
	PostalCode    string   `json:"postalCode"`
	PlusFour      string   `json:"plusFour"`
	Commercial    bool     `json:"commercial"`
	Predirection  string   `json:"predirection"`
	District      string   `json:"district"`
	Street        string   `json:"street"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	TimeZone      string   `json:"timeZone"`
	UtcOffset     *float64 `json:"utcOffset"`
}

type LocationService interface {
	ExtractAndEnrichLocation(ctx context.Context, tenant, address string) (*Location, error)
}

type locationService struct {
	log      logger.Logger
	services *Services
	aiModel  ai.AiModel
}

func NewLocationService(log logger.Logger, services *Services) LocationService {
	return &locationService{
		log:      log,
		services: services,
		aiModel: ai.NewAiModel(ai.AnthropicModelType, aiConfig.Config{
			Anthropic: aiConfig.AiModelConfigAnthropic{
				ApiPath: services.GlobalConfig.InternalServices.AiApiConfig.Url,
				ApiKey:  services.GlobalConfig.InternalServices.AiApiConfig.ApiKey,
				Model:   constants.AnthropicApiModel,
			},
		}),
	}
}

func (s *locationService) ExtractAndEnrichLocation(ctx context.Context, tenant, address string) (*Location, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationEventHandler.ExtractAndEnrichLocation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("address", address))

	if strings.TrimSpace(address) == "" {
		return nil, errors.New("address is empty")
	}

	// Step 1: Check if mapping exists
	locationMapping, err := s.services.PostgresRepositories.AiLocationMappingRepository.GetLatestLocationMappingByInput(ctx, address)
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
	prompt := fmt.Sprintf(s.services.GlobalConfig.InternalServices.AiApiConfig.AnthropicPrompts.LocationEnrichmentPrompt, address)
	promptLog := postgresEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      common.GetAppSourceFromContext(ctx),
		Provider:       constants.Anthropic,
		Model:          constants.AnthropicApiModel,
		PromptType:     constants.PromptTypeExtractLocationValue,
		Tenant:         &tenant,
		PromptTemplate: &s.services.GlobalConfig.InternalServices.AiApiConfig.AnthropicPrompts.LocationEnrichmentPrompt,
		Prompt:         prompt,
	}
	promptStoreLogId, err := s.services.PostgresRepositories.AiPromptLogRepository.Store(promptLog)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error storing prompt log: %v", err)
	}

	aiResult, err := s.aiModel.Inference(ctx, prompt)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get AI response"))
		s.log.Errorf("Error invoking AI: %s", err.Error())
		storeErr := s.services.PostgresRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, errors.Wrap(storeErr, "failed to update prompt log with error"))
			s.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
		return nil, err
	} else {
		storeErr := s.services.PostgresRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId, aiResult)
		if storeErr != nil {
			tracing.TraceErr(span, errors.Wrap(storeErr, "failed to update prompt log with ai response"))
			s.log.Errorf("Error updating prompt log with ai response: %v", storeErr)
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
	err = s.services.PostgresRepositories.AiLocationMappingRepository.AddLocationMapping(ctx, *locationMapping)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to store location mapping"))
		s.log.Errorf("Error storing location mapping: %v", err)
	}

	return &location, nil
}
