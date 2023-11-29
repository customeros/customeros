package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func NewAnthropicClient(cfg *AiModelConfigAnthropic, logger logger.Logger) *AnthropicClient {
	return &AnthropicClient{
		cfg:    cfg,
		logger: logger,
	}
}

type AnthropicClient struct {
	cfg    *AiModelConfigAnthropic
	logger logger.Logger
}

func InvokeAnthropic(ctx context.Context, cfg *AiModelConfigAnthropic, logger logger.Logger, prompt string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.invokeAnthropic")
	defer span.Finish()

	reqBody := map[string]interface{}{
		"prompt": prompt,
		"model":  "claude-2",
	}

	jsonBody, _ := json.Marshal(reqBody)
	reqReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", cfg.ApiPath+"/ask", reqReader)
	if err != nil {
		tracing.TraceErr(span, err)
		logger.Errorf("Error creating request: %v", err.Error())
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set(constants.ApiKeyHeader, cfg.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	response := strings.TrimSpace(data["completion"])
	span.LogFields(log.String("anthropicResponse", response))

	return response, nil
}
