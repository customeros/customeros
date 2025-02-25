package embedding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SegmentRequest struct {
	Content        string `json:"content"`
	Tokenizer      string `json:"tokenizer"`
	ReturnChunks   bool   `json:"return_chunks"`
	MaxChunkLength int    `json:"max_chunk_length"`
}

type SegmentResponse struct {
	NumTokens int      `json:"num_tokens"`
	Tokenizer string   `json:"tokenizer"`
	NumChunks int      `json:"num_chunks"`
	Chunks    []string `json:"chunks"`
}

const SegmentURL = "https://api.jina.ai/v1/segment"

func (s *embeddingService) Segment(ctx context.Context, input string, maxLength int) (*SegmentResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "embeddingService.Segment")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	body := SegmentRequest{
		Content:        input,
		Tokenizer:      "cl100k_base",
		ReturnChunks:   true,
		MaxChunkLength: 1024,
	}

	resp, err := s.postRequest(ctx, body, SegmentURL)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	parsedResp, err := s.parseSegments(ctx, resp)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return parsedResp, nil
}

func (s *embeddingService) parseSegments(ctx context.Context, jsonStr string) (*SegmentResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "embeddingService.parseSegments")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var response SegmentResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chunking response: %w", err)
	}

	// Validate that the number of chunks matches the chunks array length
	if response.NumChunks != len(response.Chunks) {
		return nil, fmt.Errorf("inconsistent chunk count: stated %d but got %d chunks",
			response.NumChunks, len(response.Chunks))
	}

	return &response, nil
}
