package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type embeddingService struct {
	config *config.JinaConfig
}

func NewEmbeddingService(config *config.JinaConfig) interfaces.EmbeddingService {
	return &embeddingService{
		config: config,
	}
}

var ErrPaymentRequired = errors.New("Jina balance requires topup")

func (s *embeddingService) postRequest(ctx context.Context, body any, url string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "embeddingService.postRequest")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	authorization := "Bearer " + s.config.ApiKey

	jsonReqBody, err := json.Marshal(body)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case 402:
			return "", ErrPaymentRequired

		default:
			err = fmt.Errorf("error code: %d", resp.StatusCode)
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	responseStr := string(bodyBytes)
	return responseStr, nil
}
