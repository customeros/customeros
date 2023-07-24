package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/anthorpic-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/anthorpic-api/dto"
	"net/http"
)

type anthropicService struct {
	cfg *config.Config
}

type AnthropicService interface {
	QueryAnthropic(request dto.AnthropicApiRequest) dto.AnthropicApiResponse
}

func (s *anthropicService) QueryAnthropic(request dto.AnthropicApiRequest) dto.AnthropicApiResponse {
	reqBody := map[string]interface{}{
		"prompt":               "Human: " + request.Prompt + "\nAssistant:",
		"model":                request.Model,
		"max_tokens_to_sample": request.MaxTokensToSample,
		"temperature":          request.Temperature,
		"stop_sequences":       request.StopSequences,
	}

	jsonBody, _ := json.Marshal(reqBody)
	reqReader := bytes.NewReader(jsonBody)

	req, _ := http.NewRequest("POST", s.cfg.Anthropic.ApiPath, reqReader)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", s.cfg.Anthropic.ApiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	var data dto.AnthropicApiResponse
	json.NewDecoder(resp.Body).Decode(&data)

	return data
}

func NewAnthropicService(cfg *config.Config) AnthropicService {
	return &anthropicService{
		cfg: cfg,
	}
}
