package embedding

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

func (s *embeddingService) EmbedWebpage(ctx context.Context, webpage postgres_entity.GlobalOrganizationWebpages) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "embeddingService.EmbedWebpage")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// ensure index exists in OpenSearch
	index := fmt.Sprintf("%s-%s-%s", "embeddings", "webpage", webpage.PrimaryDomain)
	err := s.opensearchService.EmbeddingsIndexCheck(ctx, index)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// clean webpage content
	cleanContent, err := s.cleanWebpage(ctx, webpage.Content)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if cleanContent == nil {
		err := errors.New("Unable to clean webpage content")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// segment content
	segments, err := s.Segment(ctx, *cleanContent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if segments == nil {
		err := errors.New("Unable to segment content")
		tracing.TraceErr(span, err)
		return nil, err
	}
	if len(segments.Chunks) == 0 {
		err := errors.New("Segment chunks are empty")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var created []string
	for _, segment := range segments.Chunks {
		// build embedding record for each segment
		record := s.newEmbeddingRecord("WEBPAGE", segment, &webpage.UpdatedAt)
		record.SourceUrl = webpage.Url

		// embed webpage content
		embeddings, err := s.GetEmbedding(ctx, segment, enum.EmbeddingPassageRetrieval)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}
		if embeddings == nil {
			continue
		}
		record.Vector = embeddings
		record.EmbeddedAt = utils.Now()

		// generate summary
		summary, err := s.generateSummary(ctx, segment)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}

		// embed summary
		summaryEmbeddings, err := s.GetEmbedding(ctx, *summary, enum.EmbeddingPassageRetrieval)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}
		if summaryEmbeddings == nil {
			continue
		}

		record.Summary = *summary
		record.SummaryVector = summaryEmbeddings

		// generate questions for content
		questions, err := s.generateQuestionsForWebsiteContent(ctx, segment)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}

		// embed questions
		for _, question := range questions {
			embeddings, err := s.GetEmbedding(ctx, question, enum.EmbeddingQuery)
			if err != nil {
				tracing.TraceErr(span, err)
				continue
			}
			if embeddings == nil {
				continue
			}
			q := dto.Question{
				Text:   question,
				Vector: embeddings,
			}
			record.Questions = append(record.Questions, q)
		}

		// generate tags for content

		// store each segment in opensearch
		err = s.opensearchService.UpsertDocument(ctx, index, &record.ID, record)
		if err != nil {
			tracing.TraceErr(span, err)
			continue
		}
		created = append(created, record.ID)
	}

	return created, nil
}

func (s *embeddingService) cleanWebpage(ctx context.Context, content string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "embeddingService.cleanWebpage")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	systemPrompt := "I'm going to give you scraped markdown content from a B2B company's website.  Your job is to remove all links, cookie warnings, menus, headers, footers, etc and return only the core page content. Remove all extra whitespace. Ensure you do not miss any content. Return only the exact text on the page, nothing else. Do not add your own comments or preamble."

	temperature := float32(0.1)
	maxOutput := int32(8000)

	answer, err := s.aiService.AskAI(ctx, interfaces.AskAIRequest{
		Model:            enum.AIModelLlama8B,
		SystemPrompt:     &systemPrompt,
		Prompt:           &content,
		ModelTemperature: &temperature,
		MaxOutputTokens:  &maxOutput,
		OutputFormat:     enum.AIOutputText,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	noWhitespace := strings.Join(strings.Fields(*answer), " ")

	return &noWhitespace, nil
}

func (s *embeddingService) generateQuestionsForWebsiteContent(ctx context.Context, content string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "generateQuestionsForContent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	systemPrompt := "I'm going to give you a content segment from a B2B company's website.  I want you to give me back a list of questions you believe this content block answers fully.  Return only questions, no answers.  Format your response as a json array."

	temperature := float32(1.0)
	maxOutput := int32(1024)

	answer, err := s.aiService.AskAI(ctx, interfaces.AskAIRequest{
		Model:            enum.AIModelGemini,
		SystemPrompt:     &systemPrompt,
		Prompt:           &content,
		ModelTemperature: &temperature,
		MaxOutputTokens:  &maxOutput,
		OutputFormat:     enum.AIOutputJson,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var questions []string

	err = json.Unmarshal([]byte(*answer), &questions)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return questions, nil
}

func (s *embeddingService) generateSummary(ctx context.Context, content string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "embeddingService.generateSummary")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	systemPrompt := "Please provide a short summary for the content provided below."
	temperature := float32(1.0)
	maxOutput := int32(1024)

	answer, err := s.aiService.AskAI(ctx, interfaces.AskAIRequest{
		Model:            enum.AIModelGeminiLite,
		SystemPrompt:     &systemPrompt,
		Prompt:           &content,
		ModelTemperature: &temperature,
		MaxOutputTokens:  &maxOutput,
		OutputFormat:     enum.AIOutputText,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return answer, nil
}

func (s *embeddingService) generateTags(ctx context.Context, content string) {
}
