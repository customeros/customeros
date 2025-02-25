package embedding

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type EmbeddingTask string

const (
	Classification   EmbeddingTask = "classification"
	Generic          EmbeddingTask = "text-matching"
	PassageRetrieval EmbeddingTask = "retrieval.passage"
	Query            EmbeddingTask = "retrieval.query"
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

type EmbeddingRecord struct {
	ID               string     `json:"id"`
	Vector           []float64  `json:"vector"`
	ContentID        string     `json:"contentId"`
	ContentSourceURL string     `json:"contentSourceUrl"`
	ContentType      string     `json:"contentType"`
	Content          string     `json:"content"`
	EmbeddingModel   string     `json:"embeddingModel"`
	ContentCreatedAt time.Time  `json:"contentCreatedAt"`
	EmbeddedAt       time.Time  `json:"embeddedAt"`
	Tags             []Tag      `json:"tags"`
	Summary          string     `json:"summary"`
	SummaryVector    []float64  `json:"summaryVector"`
	Questions        []Question `json:"questions"`
}

type Tag struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Question struct {
	Text   string    `json:"text"`
	Vector []float64 `json:"vector"`
}

const (
	EmbeddingModel = enum.AIModelJinaEmbeddings
	EmbeddingURL   = "https://api.jina.ai/v1/embeddings"
)

func (s *embeddingService) BuildEmbeddingRecord(ctx context.Context, content string) ([]*EmbeddingResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "embeddingService.BuildEmbeddingRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	return nil, nil
}

func (s *embeddingService) embedd(ctx context.Context, task EmbeddingTask, input []string) (*EmbeddingResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "embeddingService.embedd")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

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
		tracing.TraceErr(span, err)
		return nil, err
	}

	parsedResponse, err := s.parseEmbeddingResponse(ctx, resp)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return parsedResponse, nil
}

func (s *embeddingService) getEmbeddingByIndex(ctx context.Context, response *EmbeddingResponse, index int) ([]float64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getEmbeddingByIndex")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if index < 0 || index >= len(response.Data) {
		return nil, fmt.Errorf("embedding index %d out of range (0-%d)", index, len(response.Data)-1)
	}
	return response.Data[index].Embedding, nil
}

func (s *embeddingService) parseEmbeddingResponse(ctx context.Context, response string) (*EmbeddingResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "embeddingService.parseEmbeddingResponse")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var data EmbeddingResponse
	err := json.Unmarshal([]byte(response), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse embedding response: %w", err)
	}
	return &data, nil
}
