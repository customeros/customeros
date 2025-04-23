package ipdata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
)

func (s *IPDataService) AskIPData(ctx context.Context, ipAddress string) (*IPDataResponseBody, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ipDataService.AskIPData")
	defer spans.Finish()

	// validate if IPData is configured
	if s.config.ApiKey == "" || s.config.ApiUrl == "" {
		err := errors.New("IPData is not configured")
		spans.TraceError(err)
		return nil, err
	}

	// Create HTTP client
	client := &http.Client{}

	// Create IPData request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?api-key=%s", s.config.ApiUrl, ipAddress, s.config.ApiKey), nil)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to create GET request for IPData"))
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		wrappedErr := errors.Wrap(err, "failed to perform GET request for IPData")
		spans.TraceError(wrappedErr)
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "failed to read response body"))
		return nil, err
	}

	// knownBadResponse := false
	// if resp.StatusCode != http.StatusOK {
	// 	if resp.StatusCode == http.StatusBadRequest {
	// 		for _, msg := range knowIpDataBadResponseMessages {
	// 			if strings.Contains(string(responseBody), msg) {
	// 				knownBadResponse = true
	// 				break
	// 			}
	// 		}
	// 	}
	// 	if !knownBadResponse {
	// 		spans.LogKV("response.body", string(responseBody))
	// 		spans.TraceError(errors.Errorf("IPData returned status code %d", resp.StatusCode))
	// 	}
	// }

	// Parse the JSON request body
	var ipDataResponseBody IPDataResponseBody
	if err = json.Unmarshal(responseBody, &ipDataResponseBody); err != nil {
		spans.TraceError(errors.Wrap(err, "failed to unmarshal response body"))
		return nil, err
	}
	ipDataResponseBody.StatusCode = resp.StatusCode

	return &ipDataResponseBody, nil
}
