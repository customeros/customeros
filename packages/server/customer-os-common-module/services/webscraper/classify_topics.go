package webscraper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type WebpageTopics struct {
	Topics []string `json:"topics"`
}

func (s *webscraperService) ClassifyWebpageTopics(ctx context.Context, url string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "webscraperService.ClassifyWebpageTopics")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	webpage, err := s.postgresRepositories.ScrapedWebpageRepository.GetWebpage(ctx, url, 365)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	systemPrompt := `I will provide you with the scraped content of a webpage and a brief description of the company who owns it. Your job is to analyze the website content and give me up to a maximum of 5 topics or categories that most accurately describe the page. Please respond with valid json in this exact format: 
{
  "topics": [
    "topic 1",
    "topic 2", 
    "topic 3"
  ]
}`

	globalOrg, err := s.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, webpage.PrimaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Company name: %s", globalOrg.Name))
	prompt.WriteString(fmt.Sprintf("Company description: %s", globalOrg.Description))
	prompt.WriteString(fmt.Sprintf("Webpage url: %s", webpage.Url))
	prompt.WriteString("--- Webpage content --- ")
	prompt.WriteString(webpage.Content)
	promptStr := prompt.String()

	temperature := float32(0.2)
	maxOutputTokens := int32(100)
	answer, err := s.aiService.AskAI(ctx, interfaces.AskAIRequest{
		Model:            enum.AIModelLlama8B,
		SystemPrompt:     &systemPrompt,
		Prompt:           &promptStr,
		ModelTemperature: &temperature,
		MaxOutputTokens:  &maxOutputTokens,
		OutputFormat:     enum.AIOutputJson,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if answer == nil {
		return nil, nil
	}

	topics, err := s.parseTopics(*answer)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	err = s.postgresRepositories.ScrapedWebpageRepository.SetWebpageTopics(ctx, url, topics)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return topics, nil
}

func (s *webscraperService) parseTopics(answer string) ([]string, error) {
	var response WebpageTopics

	// Unmarshal the JSON string into the struct
	err := json.Unmarshal([]byte(answer), &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return response.Topics, nil
}
