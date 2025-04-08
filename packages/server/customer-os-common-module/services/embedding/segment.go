package embedding

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type SegmentRequest struct {
	Content        string `json:"content"`
	Tokenizer      string `json:"tokenizer"`
	ReturnChunks   bool   `json:"return_chunks"`
	MaxChunkLength int    `json:"max_chunk_length"`
}

const SegmentURL = "https://api.jina.ai/v1/segment"

func (s *embeddingService) Segment(ctx context.Context, input string) (*interfaces.ContentSegments, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "embeddingService.Segment")
	defer spans.Finish()

	body := SegmentRequest{
		Content:        input,
		Tokenizer:      "cl100k_base",
		ReturnChunks:   true,
		MaxChunkLength: 1024,
	}

	resp, err := s.postRequest(ctx, body, SegmentURL)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	parsedResp, err := s.parseSegments(ctx, resp)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return parsedResp, nil
}

func (s *embeddingService) parseSegments(ctx context.Context, jsonStr string) (*interfaces.ContentSegments, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "embeddingService.parseSegments")
	defer spans.Finish()

	var response interfaces.ContentSegments
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
