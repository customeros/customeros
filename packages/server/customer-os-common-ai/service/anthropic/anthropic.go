package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
)

const ApiKeyHeader = "X-Openline-API-KEY"

func NewAnthropicClient(cfg *AiModelConfigAnthropic) *AnthropicClient {
	return &AnthropicClient{
		cfg: cfg,
	}
}

type AnthropicClient struct {
	cfg *AiModelConfigAnthropic
}

func InvokeAnthropic(ctx context.Context, cfg *AiModelConfigAnthropic, prompt string) (string, error) {
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
		opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, err)
		logrus.Errorf("Error creating request: %v", err.Error())
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set(ApiKeyHeader, cfg.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, err)
		logrus.Errorf("Error executing request: %v", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	response := strings.TrimSpace(data["completion"])
	span.LogFields(log.String("anthropicResponse", response))
	logrus.Info("Completed executing Anthropic request")

	return response, nil
}
