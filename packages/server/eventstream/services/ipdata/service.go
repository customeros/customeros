package ipdata

import "context"

type IPDataService interface {
	AskIPData(ctx context.Context, ipAddress string) *IPDataResponseBody
}

func (s *verifyService) askIpData(ctx context.Context, ip string) (*postgresentity.IPDataResponseBody, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IpIntelligenceService.askIpData")
	defer spans.Finish()

	// validate if IPData is configured
	if s.cfg.External.IpDataConfig.ApiKey == "" || s.cfg.External.IpDataConfig.ApiUrl == "" {
		err := errors.New("IPData is not configured")
		spans.TraceError(err)
		s.log.Errorf("IPData is not configured")
		return nil, err
	}

	// Create HTTP client
	client := &http.Client{}

	// Create IPData request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?api-key=%s", s.cfg.External.IpDataConfig.ApiUrl, ip, s.cfg.External.IpDataConfig.ApiKey), nil)
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

	knownBadResponse := false
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			for _, msg := range knowIpDataBadResponseMessages {
				if strings.Contains(string(responseBody), msg) {
					knownBadResponse = true
					break
				}
			}
		}
		if !knownBadResponse {
			spans.LogKV("response.body", string(responseBody))
			spans.TraceError(errors.Errorf("IPData returned status code %d", resp.StatusCode))
		}
	}

	// Parse the JSON request body
	var ipDataResponseBody postgresentity.IPDataResponseBody
	if err = json.Unmarshal(responseBody, &ipDataResponseBody); err != nil {
		spans.TraceError(errors.Wrap(err, "failed to unmarshal response body"))
		return nil, err
	}
	ipDataResponseBody.StatusCode = resp.StatusCode

	return &ipDataResponseBody, nil
}
