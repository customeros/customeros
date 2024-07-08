package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/dto"
	"io"
	"net/http"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-ai/config"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
)

const ApiKeyHeader = "X-Openline-API-KEY"

func NewAnthropicClient(cfg *config.AiModelConfigAnthropic) *AnthropicClient {
	return &AnthropicClient{
		cfg: cfg,
	}
}

type AnthropicClient struct {
	cfg *config.AiModelConfigAnthropic
}

func InvokeAnthropic(ctx context.Context, cfg *config.AiModelConfigAnthropic, prompt string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.invokeAnthropic")
	defer span.Finish()
	span.LogFields(log.String("anthropicPrompt", prompt))

	reqBody := map[string]interface{}{
		"prompt": prompt,
		"model":  cfg.Model,
	}

	jsonBody, _ := json.Marshal(reqBody)
	reqReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", cfg.ApiPath+"/ask-anthropic", reqReader)
	if err != nil {
		opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, err)
		logrus.Errorf("Error creating request: %v", err.Error())
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set(ApiKeyHeader, cfg.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logrus.Errorf("Error executing request: %v", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var errorResponse struct {
			Type  string `json:"type"`
			Error struct {
				Type    string `json:"type"`
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.Unmarshal(body, &errorResponse); err != nil {
			// If we can't parse the error response, return the raw body
			return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("%s: %s", errorResponse.Error.Type, errorResponse.Error.Message)
	}

	var data dto.AnthropicApiResponse
	json.NewDecoder(resp.Body).Decode(&data)
	response := strings.TrimSpace(data.Content[0].Text)
	span.LogFields(log.String("anthropicResponse", response))
	logrus.Info("Completed executing Anthropic request")

	return response, nil
}
