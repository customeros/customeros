package service

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/openai-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/openai-api/dto"
	"net/http"
	"strings"
)

type openAiService struct {
	cfg *config.Config
}

type OpenAiService interface {
	QueryOpenAi(request dto.OpenAiApiRequest) dto.OpenAiApiResponse
}

func (s *openAiService) QueryOpenAi(request dto.OpenAiApiRequest) dto.OpenAiApiResponse {
	requestData := map[string]interface{}{}
	requestData["model"] = request.Model
	requestData["max_tokens"] = request.MaxTokensToSample
	requestData["temperature"] = request.Temperature
	requestData["messages"] = []interface{}{}
	requestData["messages"] = append(requestData["messages"].([]interface{}), map[string]interface{}{})
	requestData["messages"].([]interface{})[0].(map[string]interface{})["role"] = "user"
	requestData["messages"].([]interface{})[0].(map[string]interface{})["content"] = request.Prompt

	requestBody, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", s.cfg.OpenAi.ApiPath, strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.OpenAi.ApiKey)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making the API request:", err)
		return dto.OpenAiApiResponse{
			Error: &dto.OpenAiApiErrorResponse{
				Type:    "REQUEST_ERROR",
				Message: err.Error(),
			},
		}
	}
	if response.StatusCode != 200 {
		fmt.Println("Error making the API request:", response.Status)
		return dto.OpenAiApiResponse{
			Error: &dto.OpenAiApiErrorResponse{
				Type:    "REQUEST_ERROR",
				Message: "HTTP status code: " + response.Status,
			},
		}
	}
	defer response.Body.Close()

	var result dto.OpenAiApiResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error parsing the API response:", err)
		return dto.OpenAiApiResponse{
			Error: &dto.OpenAiApiErrorResponse{
				Type:    "REQUEST_ERROR",
				Message: err.Error(),
			},
		}
	}

	return result
}

func NewOpenAiService(cfg *config.Config) OpenAiService {
	return &openAiService{
		cfg: cfg,
	}
}
