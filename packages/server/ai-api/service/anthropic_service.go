package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/dto"
)

type anthropicService struct {
	cfg *config.Config
}

type AnthropicService interface {
	QueryAnthropic(request dto.AnthropicApiRequest) (dto.AnthropicApiResponse, error)
}

func NewAnthropicService(cfg *config.Config) AnthropicService {
	return &anthropicService{
		cfg: cfg,
	}
}

func (s *anthropicService) QueryAnthropic(request dto.AnthropicApiRequest) (dto.AnthropicApiResponse, error) {
	reqBody := map[string]interface{}{
		"model": request.Model,
		"messages": []map[string]string{
			{"role": "user", "content": request.Prompt},
		},
		"max_tokens": request.MaxTokensToSample,
	}

	// Only include temperature if it's not zero
	if request.Temperature != 0 {
		reqBody["temperature"] = request.Temperature
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return dto.AnthropicApiResponse{}, fmt.Errorf("error marshaling request body: %v", err)
	}
	reqReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", s.cfg.Anthropic.ApiPath, reqReader)
	if err != nil {
		return dto.AnthropicApiResponse{}, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.cfg.Anthropic.ApiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return dto.AnthropicApiResponse{}, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.AnthropicApiResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return dto.AnthropicApiResponse{}, fmt.Errorf("error response: %s", string(body))
	}

	var data dto.AnthropicApiResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return dto.AnthropicApiResponse{}, fmt.Errorf("error decoding response: %v", err)
	}

	return data, nil
}
