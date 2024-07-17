package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"io"
	"net/http"
)

type EnrichDetailsTrackingService interface {
	ShouldIdentifyTrackingRecords(ctx context.Context) error
	IdentifyTrackingRecords(ctx context.Context) error
}

type enrichDetailsTrackingService struct {
	cfg      *config.Config
	services *Services
}

func NewEnrichDetailsTrackingService(cfg *config.Config, services *Services) EnrichDetailsTrackingService {
	return &enrichDetailsTrackingService{
		cfg:      cfg,
		services: services,
	}
}

func (s *enrichDetailsTrackingService) ShouldIdentifyTrackingRecords(ctx context.Context) error {
	span, ctx := tracing.StartTracerSpan(ctx, "EnrichDetailsTrackingService.ShouldIdentifyTrackingRecords")
	defer span.Finish()

	needsIdentificationRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetForPrefilterBeforeIdentification(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, record := range needsIdentificationRecords {
		shouldIdentify, err := s.processShouldIdentifyRecord(ctx, record.IP)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if shouldIdentify != nil {
			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.ShouldIdentify(ctx, record.ID, shouldIdentify)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}
	}

	return nil
}

func (s *enrichDetailsTrackingService) IdentifyTrackingRecords(ctx context.Context) error {
	span, ctx := tracing.StartTracerSpan(ctx, "EnrichDetailsTrackingService.IdentifyTrackingRecords")
	defer span.Finish()

	notIdentifiedTrackingRecords, err := s.services.CommonServices.PostgresRepositories.TrackingRepository.GetReadyForIdentification(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, record := range notIdentifiedTrackingRecords {
		enrichDetailsTrackingId, err := s.processNotIdentifiedRecord(ctx, record.IP)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if enrichDetailsTrackingId != nil {
			err = s.services.CommonServices.PostgresRepositories.TrackingRepository.MarkAsIdentified(ctx, record.ID)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
		}

	}

	return nil
}

func (s *enrichDetailsTrackingService) processShouldIdentifyRecord(ctx context.Context, ip string) (*bool, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "EnrichDetailsTrackingService.processShouldIdentifyRecord")
	defer span.Finish()

	ipDataByIp, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.GetByIP(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get better contact details: %v", err)
	}

	if ipDataByIp == nil {
		ipDataByIp, err = s.askAndStoreIPData(ctx, ip)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

	}

	if ipDataByIp == nil {
		tracing.TraceErr(span, errors.New("ipdata record is nil"))
	}

	return &ipDataByIp.ShouldIdentify, nil
}

func (s *enrichDetailsTrackingService) processNotIdentifiedRecord(ctx context.Context, ip string) (*string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "EnrichDetailsTrackingService.processNotIdentifiedRecord")
	defer span.Finish()

	snitcherByIp, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get better contact details: %v", err)
	}

	if snitcherByIp == nil {
		snitcherByIp, err = s.askAndStoreSnitcherData(ctx, ip)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	if snitcherByIp == nil {
		tracing.TraceErr(span, errors.New("snitcher record is nil"))
	}

	return &snitcherByIp.ID, nil
}

func (s *enrichDetailsTrackingService) askAndStoreIPData(ctx context.Context, ip string) (*entity.EnrichDetailsPreFilterTracking, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "EnrichDetailsTrackingService.askAndStoreIPData")
	defer span.Finish()

	// Create HTTP client
	client := &http.Client{}

	// Create POST request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?api-key=%s", s.cfg.IPDataApi.Url, ip, s.cfg.IPDataApi.ApiKey), nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	//Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON request body
	var ipDataResponseBody entity.IPDataResponseBody
	if err = json.Unmarshal(responseBody, &ipDataResponseBody); err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	shouldIdentify := true

	if ipDataResponseBody.Ip == "" {
		shouldIdentify = false
	}

	if ipDataResponseBody.Carrier != nil {
		shouldIdentify = false
	}

	if ipDataResponseBody.Threat.IsTor ||
		ipDataResponseBody.Threat.IsIcloudRelay ||
		ipDataResponseBody.Threat.IsProxy ||
		ipDataResponseBody.Threat.IsDatacenter ||
		ipDataResponseBody.Threat.IsAnonymous ||
		ipDataResponseBody.Threat.IsKnownAttacker ||
		ipDataResponseBody.Threat.IsKnownAbuser ||
		ipDataResponseBody.Threat.IsThreat ||
		ipDataResponseBody.Threat.IsBogon {
		shouldIdentify = false
	}

	marshal, err := json.Marshal(ipDataResponseBody)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to marshal response body: %v", err)
	}

	// Store response
	err = s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.RegisterRequest(ctx, entity.EnrichDetailsPreFilterTracking{
		CreatedAt:      utils.Now(),
		IP:             ip,
		ShouldIdentify: shouldIdentify,
		Response:       string(marshal),
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to store response: %v", err)
	}

	byIP, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsPrefilterTrackingRepository.GetByIP(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get stored response: %v", err)
	}

	return byIP, nil
}

func (s *enrichDetailsTrackingService) askAndStoreSnitcherData(ctx context.Context, ip string) (*entity.EnrichDetailsTracking, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "EnrichDetailsTrackingService.askAndStoreSnitcherData")
	defer span.Finish()

	// Create HTTP client
	client := &http.Client{}

	// Create POST request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/company/find?ip=%s", s.cfg.SnitcherApi.Url, ip), nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.SnitcherApi.ApiKey)

	//Perform the request
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to perform POST request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON request body
	var snitherResponse entity.SnitcherResponseBody
	if err = json.Unmarshal(responseBody, &snitherResponse); err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	var companyName, companyDomain *string

	if snitherResponse.Company != nil && snitherResponse.Company.Name != "" {
		companyName = &snitherResponse.Company.Name
	}

	if snitherResponse.Company != nil && snitherResponse.Company.Domain != "" {
		companyDomain = &snitherResponse.Company.Domain
	}

	// Store response
	err = s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.RegisterRequest(ctx, entity.EnrichDetailsTracking{
		CreatedAt:     utils.Now(),
		IP:            ip,
		CompanyName:   companyName,
		CompanyDomain: companyDomain,
		Response:      string(responseBody),
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to store response: %v", err)
	}

	byIP, err := s.services.CommonServices.PostgresRepositories.EnrichDetailsTrackingRepository.GetByIP(ctx, ip)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get stored response: %v", err)
	}

	return byIP, nil
}
