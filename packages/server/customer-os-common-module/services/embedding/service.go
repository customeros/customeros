package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"io"
	"net/http"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type embeddingService struct {
	config                *config.JinaConfig
	aiService             interfaces.AIService
	opensearchService     interfaces.OpensearchService
	warehouseRepositories *postgres_repository.WarehouseRepositories
}

func NewEmbeddingService(config *config.JinaConfig, aiService interfaces.AIService, opensearch interfaces.OpensearchService, warehouse *postgres_repository.WarehouseRepositories) interfaces.EmbeddingService {
	return &embeddingService{
		config:                config,
		aiService:             aiService,
		opensearchService:     opensearch,
		warehouseRepositories: warehouse,
	}
}

var ErrPaymentRequired = errors.New("Jina balance requires topup")

func (s *embeddingService) postRequest(ctx context.Context, body any, url string) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "embeddingService.postRequest")
	defer spans.Finish()

	authorization := "Bearer " + s.config.ApiKey

	jsonReqBody, err := json.Marshal(body)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)

	// Create HTTP client
	clientTimeout := 15 * time.Second
	httpClient := clients.NewLoggingClient(s.warehouseRepositories.APICallLogRepository, enum.VendorJina, &clientTimeout)

	resp, err := httpClient.Do(req)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case 402:
			return "", ErrPaymentRequired

		default:
			err = fmt.Errorf("error code: %d", resp.StatusCode)
			spans.TraceError(err)
			return "", err
		}
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		spans.TraceError(err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	responseStr := string(bodyBytes)
	return responseStr, nil
}
