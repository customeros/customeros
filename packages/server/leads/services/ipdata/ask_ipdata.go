package ipdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/leads/internal/clients"
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

var knownBadResponseMessages = []string{"is a reserved IP address"}

func (s *IPDataService) AskIPData(ctx context.Context, ipAddress string) *pb.IPAddressVerifyResponse {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ipDataService.AskIPData")
	defer spans.Finish()

	response := &pb.IPAddressVerifyResponse{}

	// validate if IPData is configured
	if s.config.ApiKey == "" || s.config.ApiUrl == "" {
		err := errors.New("IPData is not configured")
		spans.TraceError(err)
		response.ErrorMessage = err.Error()
		return response
	}

	// Create HTTP client
	client := clients.NewLoggingClient(s.repositories.APICallLogRepository, enum.VendorIPData)

	// Create IPData request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?api-key=%s", s.config.ApiUrl, ipAddress, s.config.ApiKey), nil)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to create GET request for IPData"))
		response.ErrorMessage = err.Error()
		return response
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		wrappedErr := errors.Wrap(err, "failed to perform GET request for IPData")
		spans.TraceError(wrappedErr)
		response.ErrorMessage = wrappedErr.Error()
		return response
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusOK:
		data, err := s.handleSuccess(ctx, resp)
		if err != nil {
			spans.TraceError(err)
			response.ErrorMessage = err.Error()
			return response
		}
		return &pb.IPAddressVerifyResponse{
			IpAddress:   data.IPAddress,
			City:        data.City,
			Region:      data.Region,
			CountryCode: data.CountryCode,
			IsThreat:    data.HasThreat,
		}

	case resp.StatusCode == http.StatusBadRequest:
		err := s.handleBadRequest(ctx, resp)
		response.ErrorMessage = err.Error()
		return response

	default:
		response.ErrorMessage = "IPData unable to process request"
		return response
	}
}

func (s *IPDataService) handleBadRequest(ctx context.Context, resp *http.Response) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "IPDataService.handleBadRequest")
	defer span.Finish()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		span.TraceError(errors.Wrap(err, "failed to read response body"))
		return err
	}

	for _, msg := range knownBadResponseMessages {
		if strings.Contains(string(responseBody), msg) {
			return nil
		}
	}
	return errors.New("IP Data returned Bad Request")
}

func (s *IPDataService) handleSuccess(ctx context.Context, resp *http.Response) (*models.IPIntelligence, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "IPDataService.handleSuccess")
	defer span.Finish()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		span.TraceError(errors.Wrap(err, "failed to read response body"))
		return nil, err
	}

	// Parse the JSON request body
	var ipDataResponseBody IPDataResponseBody
	if err := json.Unmarshal(responseBody, &ipDataResponseBody); err != nil {
		span.TraceError(errors.Wrap(err, "failed to unmarshal response body"))
		return nil, err
	}
	ipDataResponseBody.StatusCode = resp.StatusCode

	// check if record exists
	record, err := s.repositories.IPIntelligence.FindByIP(ctx, ipDataResponseBody.Ip)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	newRecord := buildIPIntelligenceRecord(&ipDataResponseBody)

	if record != nil {
		updatedRecord, err := s.updateIpIntelligenceRecord(ctx, record, newRecord)
		if err != nil {
			span.TraceError(err)
			return nil, err
		}
		return updatedRecord, nil
	}

	err = s.repositories.IPIntelligence.Create(ctx, newRecord)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}
	return newRecord, nil
}

func (s *IPDataService) updateIpIntelligenceRecord(ctx context.Context, existingRecord, newRecord *models.IPIntelligence) (*models.IPIntelligence, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "IPDataService.updateIpIntelligenceRecord")
	defer span.Finish()

	updatedRecord := &models.IPIntelligence{
		ID:           existingRecord.ID,
		IPAddress:    chooseMostRecent(existingRecord.IPAddress, newRecord.IPAddress),
		Domain:       existingRecord.Domain,
		DomainSource: existingRecord.DomainSource,
		EmailAddress: existingRecord.EmailAddress,
		IsMobile:     newRecord.IsMobile,
		City:         chooseMostRecent(existingRecord.City, newRecord.City),
		Region:       chooseMostRecent(existingRecord.Region, newRecord.Region),
		CountryCode:  chooseMostRecent(existingRecord.CountryCode, newRecord.CountryCode),
		HasThreat:    newRecord.HasThreat,
		CreatedAt:    existingRecord.CreatedAt,
		UpdatedAt:    newRecord.UpdatedAt,
	}

	err := s.repositories.IPIntelligence.Update(ctx, updatedRecord)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	return updatedRecord, nil
}

func buildIPIntelligenceRecord(response *IPDataResponseBody) *models.IPIntelligence {
	isMobile := false
	if response.Carrier != nil {
		isMobile = true
	}

	return &models.IPIntelligence{
		ID:          utils.GenerateNanoIDWithPrefix("ip", 21),
		IPAddress:   response.Ip,
		IsMobile:    isMobile,
		City:        response.City,
		Region:      response.Region,
		CountryCode: response.CountryCode,
		HasThreat:   response.Threat.IsThreat,
		CreatedAt:   utils.Now(),
		UpdatedAt:   utils.NowPtr(),
	}
}

func chooseMostRecent(existing, new string) string {
	if new != "" {
		return new
	}
	return existing
}
