package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"io"
	"net/http"
	"net/url"
)

type EmailFinderService interface {
	FindEmail(ctx context.Context, firstName, lastName, domain string) (dto.FindEmailResponse, error)
}

type emailFinderService struct {
	config   *config.Config
	Services *Services
	log      logger.Logger
}

func NewEmailFinderService(config *config.Config, services *Services, log logger.Logger) EmailFinderService {
	return &emailFinderService{
		config:   config,
		Services: services,
		log:      log,
	}
}

func (s *emailFinderService) FindEmail(ctx context.Context, firstName, lastName, domain string) (dto.FindEmailResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailFinderService.FindEmail")
	defer span.Finish()
	span.LogFields(log.String("firstName", firstName), log.String("lastName", lastName), log.String("domain", domain))

	if firstName == "" || lastName == "" {
		span.LogFields(log.String("result.warn", "firstName and lastName must not be empty"))
		return dto.FindEmailResponse{}, fmt.Errorf("firstName and lastName must not be empty")
	}

	apiKey := s.config.HunterConfig.ApiKey
	url := s.config.HunterConfig.ApiPath
	params := map[string]string{
		"domain":     domain,
		"first_name": firstName,
		"last_name":  lastName,
		"api_key":    apiKey,
	}

	body, err := makeHunterHTTPRequest(url, params)
	if err != nil {
		tracing.TraceErr(span, err)
		return dto.FindEmailResponse{}, err
	}

	var result struct {
		Data struct {
			FirstName string  `json:"first_name"`
			LastName  string  `json:"last_name"`
			Email     string  `json:"email"`
			Score     float64 `json:"score"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tracing.TraceErr(span, err)
		return dto.FindEmailResponse{}, err
	}

	if result.Data.Email == "" {
		span.LogFields(log.String("result.info", "email not found"))
		return dto.FindEmailResponse{}, nil
	}

	response := dto.FindEmailResponse{
		FirstName: result.Data.FirstName,
		LastName:  result.Data.LastName,
		Email:     result.Data.Email,
		Score:     result.Data.Score,
	}
	tracing.LogObjectAsJson(span, "response", response)

	return response, nil
}

func makeHunterHTTPRequest(baseUrl string, params map[string]string) ([]byte, error) {
	// Construct request URL with parameters
	reqURL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	query := reqURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	reqURL.RawQuery = query.Encode()

	// Execute HTTP GET request
	resp, err := http.Get(reqURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
