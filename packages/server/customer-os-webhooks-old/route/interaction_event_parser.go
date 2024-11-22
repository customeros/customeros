package route

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func (h *EmailWebhookHandler) parseWebhookData(c *gin.Context, span opentracing.Span) (*model.PostmarkEmailWebhookData, error) {
	var webhookData model.PostmarkEmailWebhookData
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &webhookData); err != nil {
		tracing.LogObjectAsJson(span, "requestBody", body)
		return nil, err
	}

	return &webhookData, nil
}
