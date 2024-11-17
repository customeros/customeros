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
	"time"

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

	var response string
	var lastErr error
	const maxRetries = 4
	client := &http.Client{}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("POST", cfg.ApiPath+"/ask-anthropic", reqReader)
		if err != nil {
			opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, err)
			logrus.Errorf("Error creating request: %v", err.Error())
			return "", err
		}
		req.Header.Set("content-type", "application/json")
		req.Header.Set(ApiKeyHeader, cfg.ApiKey)

		resp, err := client.Do(req)
		if err != nil {
			opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, err)
			logrus.Errorf("Error executing request: %v", err.Error())
			return "", err
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var data dto.AnthropicApiResponse
			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				logrus.Errorf("Error decoding response: %v", err.Error())
				return "", err
			}
			response = strings.TrimSpace(data.Content[0].Text)
			span.LogFields(log.String("anthropicResponse", response))
			logrus.Info("Completed executing Anthropic request")
			return response, nil
		}

		// Handle non-OK response
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
			lastErr = fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		} else if errorResponse.Error.Type == "rate_limit_error" {
			logrus.Warnf("Rate limit hit: %s. Retrying... (%d/%d)", errorResponse.Error.Message, attempt, maxRetries)
			lastErr = fmt.Errorf("%s: %s", errorResponse.Error.Type, errorResponse.Error.Message)
			if attempt < maxRetries {
				time.Sleep(1 * time.Second) // Wait before retrying
				continue
			}
		} else {
			lastErr = fmt.Errorf("%s: %s", errorResponse.Error.Type, errorResponse.Error.Message)
		}

		// If not retrying, break
		break
	}

	logrus.Errorf("Failed after %d attempts: %v", maxRetries, lastErr)
	return "", lastErr
}
