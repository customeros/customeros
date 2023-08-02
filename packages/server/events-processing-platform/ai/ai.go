package ai

import (
	"bytes"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

func InvokeAnthropic(ctx context.Context, cfg *config.Config, logger logger.Logger, prompt string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.invokeAnthropic")
	defer span.Finish()

	reqBody := map[string]interface{}{
		"prompt": prompt,
		"model":  "claude-2",
	}

	jsonBody, _ := json.Marshal(reqBody)
	reqReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", cfg.Services.Anthropic.ApiPath+"/ask", reqReader)
	if err != nil {
		tracing.TraceErr(span, err)
		logger.Errorf("Error creating request: %v", err.Error())
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Openline-API-KEY", cfg.Services.Anthropic.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	defer resp.Body.Close()

	// Print summarized email
	var data map[string]string
	json.NewDecoder(resp.Body).Decode(&data)
	result := strings.TrimSpace(data["completion"])
	span.LogFields(log.String("result", result))

	return result, nil
}
