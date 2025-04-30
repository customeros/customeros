package snitcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/customeros/customeros/packages/server/leads/internal/clients"
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

const IDENTITY_SOURCE = "Snitcher"

func (s *SnitcherService) AskSnitcher(ctx context.Context, ip string) *pb.IPAddressIdentifyResponse {
	spans, ctx := telemetry.StartServiceSpan(ctx, "snitcherService.AskSnitcher")
	defer spans.Finish()

	spans.LogKV("ip", ip)

	response := &pb.IPAddressIdentifyResponse{}

	// validate if snitcher is configured
	if s.config.ApiKey == "" || s.config.Url == "" {
		err := fmt.Errorf("snitcher is not configured")
		spans.TraceError(err)
		response.ErrorMessage = err.Error()
		return response
	}

	// Create HTTP client with timeout
	client := clients.NewLoggingClient(s.repositories.APICallLogRepository, enum.VendorSnitcher)

	// Create POST request with context
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/company/find?ip=%s", s.config.Url, ip), nil)
	if err != nil {
		spans.TraceError(err)
		err := fmt.Errorf("failed to create POST request: %w", err)
		response.ErrorMessage = err.Error()
		return response
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.ApiKey)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		spans.TraceError(err)
		err := fmt.Errorf("failed to perform POST request: %w", err)
		response.ErrorMessage = err.Error()
		return response
	}
	defer resp.Body.Close()

	// Read response with size limit
	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
	if err != nil {
		spans.TraceError(err)
		err := fmt.Errorf("failed to read response body: %w", err)
		response.ErrorMessage = err.Error()
		return response
	}

	// Check status code
	spans.LogKV("response.statusCode", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		spans.LogKV("result.rawSnitcherResponse", string(responseBody))
		err := fmt.Errorf("snitcher API returned non-200 status code: %d", resp.StatusCode)
		response.ErrorMessage = err.Error()
		return response
	}

	// Validate and compact JSON
	err = validateJSON(responseBody)
	if err != nil {
		spans.TraceError(err)
		spans.LogKV("json.response.invalid", string(responseBody))
		err := fmt.Errorf("failed to process JSON response: %w", err)
		response.ErrorMessage = err.Error()
		return response
	}

	// Parse the response
	var snitcherResponse SnitcherResponse
	if err := json.Unmarshal(responseBody, &snitcherResponse); err != nil {
		spans.TraceError(err)
		spans.LogKV("json.response.parsing", string(responseBody))
		err := fmt.Errorf("failed to parse snitcher response: %w", err)
		response.ErrorMessage = err.Error()
		return response
	}

	err = s.handleSuccessResponse(ctx, ip, &snitcherResponse)
	if err != nil {
		spans.TraceError(err)
		response.ErrorMessage = err.Error()
		return response
	}

	return &pb.IPAddressIdentifyResponse{
		IpAddress: ip,
		Domain:    snitcherResponse.Domain,
	}
}

func (s *SnitcherService) handleSuccessResponse(ctx context.Context, ipAddress string, resp *SnitcherResponse) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "SnitcherService.handleSuccessResponse")
	defer span.Finish()

	// check if record exists
	record, err := s.repositories.IPIntelligence.FindByIP(ctx, ipAddress)
	if err != nil {
		span.TraceError(err)
		return err
	}

	newRecord := buildIPIntelligenceRecord(ipAddress, resp)

	if record != nil {
		err := s.updateIpIntelligenceRecord(ctx, record, newRecord)
		if err != nil {
			span.TraceError(err)
			return err
		}
		return nil
	}

	err = s.repositories.IPIntelligence.Create(ctx, newRecord)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// TODO Throw Snitcher company enrichment event to enrich global orgs table

	return nil
}

func (s *SnitcherService) updateIpIntelligenceRecord(ctx context.Context, existingRecord, newRecord *models.IPIntelligence) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "SnitcherService.updateIpIntelligenceRecord")
	defer span.Finish()

	updatedRecord := &models.IPIntelligence{
		ID:           existingRecord.ID,
		IPAddress:    chooseMostRecent(existingRecord.IPAddress, newRecord.IPAddress),
		Domain:       chooseMostRecent(existingRecord.Domain, newRecord.Domain),
		DomainSource: IDENTITY_SOURCE,
		EmailAddress: existingRecord.EmailAddress,
		IsMobile:     existingRecord.IsMobile,
		City:         existingRecord.City,
		Region:       existingRecord.Region,
		CountryCode:  existingRecord.CountryCode,
		HasThreat:    existingRecord.HasThreat,
		CreatedAt:    existingRecord.CreatedAt,
		UpdatedAt:    newRecord.UpdatedAt,
	}

	err := s.repositories.IPIntelligence.Update(ctx, updatedRecord)
	if err != nil {
		span.TraceError(err)
		return err
	}

	return nil
}

func buildIPIntelligenceRecord(ipAddress string, response *SnitcherResponse) *models.IPIntelligence {
	return &models.IPIntelligence{
		ID:           utils.GenerateNanoIDWithPrefix("ip", 21),
		IPAddress:    ipAddress,
		Domain:       response.Domain,
		DomainSource: IDENTITY_SOURCE,
		CreatedAt:    utils.Now(),
		UpdatedAt:    utils.NowPtr(),
	}
}

func validateJSON(data []byte) error {
	if !json.Valid(data) {
		return fmt.Errorf("invalid JSON")
	}
	return nil
}

func chooseMostRecent(existing, new string) string {
	if new != "" {
		return new
	}
	return existing
}
