package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type searchService struct {
	opensearch       interfaces.OpensearchService
	embeddingService interfaces.EmbeddingService
	aiService        interfaces.AIService
}

func NewSearchService(
	opensearch interfaces.OpensearchService,
	embeddingService interfaces.EmbeddingService,
	aiService interfaces.AIService,
) interfaces.SearchService {
	return &searchService{
		opensearch:       opensearch,
		embeddingService: embeddingService,
		aiService:        aiService,
	}
}

func (s *searchService) SearchWebsites(ctx context.Context, primaryDomain, query string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SearchService.HybridSearch")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Get embedding for the query
	embeddedQuery, err := s.embeddingService.GetEmbedding(ctx, query, enum.EmbeddingQuery)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", fmt.Errorf("error getting embedding for query: %w", err)
	}

	// Perform hybrid search with the generated embedding
	limit := 10
	search := interfaces.HybridSearchRequest{
		Index:         fmt.Sprintf("%s-%s-%s", "embeddings", "webpage", primaryDomain),
		Query:         query,
		EmbeddedQuery: embeddedQuery,
		ResultsLimit:  &limit,
	}
	results, err := s.opensearch.HybridSearch(ctx, search)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	context := s.buildPromptContext(ctx, results)
	return s.getAnswer(ctx, query, context)
}

func (s *searchService) buildPromptContext(ctx context.Context, results []interfaces.HybridSearchResult) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SearchService.buildPromptContext")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	var context strings.Builder
	context.WriteString("Context:\n\n")

	for i, result := range results {
		sectionHeader := fmt.Sprintf("--- Document %d (Score: %.2f) ---\n", i+1, result.Score)
		context.WriteString(sectionHeader)
		context.WriteString("Content:\n")
		context.WriteString(result.Content)
		context.WriteString("\n\n")
	}

	return context.String()
}

func (s *searchService) getAnswer(ctx context.Context, query, context string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SearchService.getAnswer")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	systemPrompt := fmt.Sprintf("I want you to answer the following question as directly as you can from the context provided.  The question I want you to answer is: %s.  Below is all the context you need to formulate your answer.  Again, please be as direct as possible in your response.  Avoid references to documents and avoid all preamble.", query)

	temperature := float32(0.1)
	answer, err := s.aiService.AskAI(ctx, interfaces.AskAIRequest{
		Model:            enum.AIModelAnthropicSonnet,
		SystemPrompt:     &systemPrompt,
		Prompt:           &context,
		ModelTemperature: &temperature,
		OutputFormat:     enum.AIOutputText,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if answer == nil {
		return "", nil
	}
	return *answer, nil
}
