package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
)

type IPIdentityService interface {
	IPIdentity(c context.Context, ip string) (*model.SnitcherResponse, error)
}

type ipIdentityService struct {
	config   *config.Config
	services *Services
	log      logger.Logger
}

func NewIPIdentityService(config *config.Config, services *Services) IPIdentityService {
	return &ipIdentityService{
		config:   config,
		services: services,
	}
}

const CACHE_LOOKBACK = 90 // days

func (s *ipIdentityService) IPIdentity(c context.Context, ip string) (*model.SnitcherResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(c, "SnitcherService.GetSnitcherData")
	defer span.Finish()

	// check to see if IP mapping data already exists
	query := entity.CacheIPIdentify{
		IPAddress: ip,
	}

	results, err := s.services.CommonServices.PostgresRepositories.CacheIPIdentifyRepository.Find(ctx, query, CACHE_LOOKBACK)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	if results != nil {
		var snitcherResponse model.SnitcherResponse
		err := json.Unmarshal([]byte(results.SnitcherData), &snitcherResponse)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling snitcher response: %w", err)
		}
		return &snitcherResponse, nil
	}

	snitcherResponse, respString, err := s.callSnitcher(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// return early if no company found
	if !snitcherResponse.CompanyFound() {
		return snitcherResponse, nil
	}

	// Store response
	snitcherData := entity.CacheIPIdentify{
		IPAddress:    ip,
		Domain:       snitcherResponse.CompanyDomain(),
		LinkedinSlug: snitcherResponse.LinkedinSlug(),
		SnitcherData: *respString,
	}

	err = s.services.CommonServices.PostgresRepositories.CacheIPIdentifyRepository.Create(ctx, snitcherData)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to store response: %v", err)
	}

	return snitcherResponse, nil
}

func (s *ipIdentityService) callSnitcher(ctx context.Context, ip string) (*model.SnitcherResponse, *string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SnitcherService.callSnitcher")
	defer span.Finish()

	// Create HTTP client
	client := &http.Client{}

	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/company/find?ip=%s", s.config.SnitcherConfig.Url, ip), nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.SnitcherConfig.ApiKey)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, responseBody); err != nil {
		return nil, nil, fmt.Errorf("error compacting json: %w", err)
	}

	responseString := buf.String()

	snitcherData, err := buildSnitcherResponse(ctx, responseBody)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, fmt.Errorf("faild to parse snitcher response: %v", err)
	}

	return snitcherData, &responseString, nil
}

func buildSnitcherResponse(ctx context.Context, responseBody []byte) (*model.SnitcherResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SnitcherService.buildSnitcherResponse")
	defer span.Finish()

	// Parse the JSON request body
	var snitcherResponse model.SnitcherResponse
	if err := json.Unmarshal(responseBody, &snitcherResponse); err != nil {
		tracing.TraceErr(span, err)
		span.LogFields(log.String("snitcherResponseBody", string(responseBody)))
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &snitcherResponse, nil
}
