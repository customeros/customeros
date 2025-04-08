package embedding

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type EmbeddingRequestBody struct {
	Model         string   `json:"model"`
	Task          string   `json:"task"`
	LateChunking  bool     `json:"late_chunking"`
	Dimensions    int      `json:"dimensions"`
	EmbeddingType string   `json:"embedding_type"`
	Input         []string `json:"input"`
}

type EmbeddingResponse struct {
	Model  string         `json:"model"`
	Object string         `json:"object"`
	Usage  Usage          `json:"usage"`
	Data   []EmbeddingObj `json:"data"`
}

type Usage struct {
	TotalTokens  int `json:"total_tokens"`
	PromptTokens int `json:"prompt_tokens"`
}

type EmbeddingObj struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

const (
	EmbeddingModel = enum.AIModelJinaEmbeddings
	EmbeddingURL   = "https://api.jina.ai/v1/embeddings"
)

func (s *embeddingService) newEmbeddingRecord(contentType enum.EmbeddingContentType, content string, contentCreatedAt *time.Time) *dto.EmbeddingRecord {
	record := &dto.EmbeddingRecord{
		ID:             utils.GenerateNanoIdWithPrefix("emb", 21),
		ContentType:    contentType.String(),
		Content:        content,
		EmbeddingModel: EmbeddingModel.String(),
	}
	if contentCreatedAt == nil {
		record.ContentCreatedAt = utils.Now()
	} else {
		record.ContentCreatedAt = *contentCreatedAt
	}
	return record
}

func (s *embeddingService) GetEmbedding(ctx context.Context, content string, task enum.EmbeddingTask) ([]float64, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "embeddingService.embedWebpageSegment")
	defer spans.Finish()

	embedding, err := s.embed(ctx, task, []string{content})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if embedding == nil {
		err := errors.New("unable to embed content")
		return nil, err
	}
	if len(embedding.Data) == 0 {
		err := errors.New("no embeddings returned")
		return nil, err
	}

	return embedding.Data[0].Embedding, nil
}

func (s *embeddingService) embed(ctx context.Context, task enum.EmbeddingTask, input []string) (*EmbeddingResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "embeddingService.embed")
	defer spans.Finish()

	if s.config.ApiKey == "" {
		return nil, errors.New("Jina API key not set")
	}

	body := EmbeddingRequestBody{
		Model:         EmbeddingModel.String(),
		Task:          string(task),
		LateChunking:  true,
		Dimensions:    1024,
		EmbeddingType: "float",
		Input:         input,
	}

	resp, err := s.postRequest(ctx, body, EmbeddingURL)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	parsedResponse, err := s.parseEmbeddingResponse(ctx, resp)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return parsedResponse, nil
}

func (s *embeddingService) getEmbeddingByIndex(ctx context.Context, response *EmbeddingResponse, index int) ([]float64, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "getEmbeddingByIndex")
	defer spans.Finish()

	if index < 0 || index >= len(response.Data) {
		return nil, fmt.Errorf("embedding index %d out of range (0-%d)", index, len(response.Data)-1)
	}
	return response.Data[index].Embedding, nil
}

func (s *embeddingService) parseEmbeddingResponse(ctx context.Context, response string) (*EmbeddingResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "embeddingService.parseEmbeddingResponse")
	defer spans.Finish()

	var data EmbeddingResponse
	err := json.Unmarshal([]byte(response), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse embedding response: %w", err)
	}
	return &data, nil
}
