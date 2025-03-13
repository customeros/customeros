package webscraper

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	"strings"

	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

func (s *webscraperService) ClassifyWebpageTopics(ctx context.Context, url string, pageContent *string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebscraperService.ClassifyWebpageTopics")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("url", url))
	tracing.LogObjectAsJson(span, "pageContent", pageContent)

	_, primaryDomain := domaincheck.PrimaryDomainCheck(utils.ExtractDomain(url))

	if pageContent == nil {
		webpage, err := s.postgresRepositories.ScrapedWebpageRepository.GetWebpage(ctx, url, 365)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if webpage == nil {
			err := errors.New("webpage doesn't exist")
			span.LogKV("url", url)
			tracing.TraceErr(span, err)
			return nil, err
		}

		pageContent = &webpage.Content
		primaryDomain = webpage.PrimaryDomain
	}

	globalOrg, err := s.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	systemPrompt := `I will provide you with the scraped content of a webpage and a brief description of the company who owns it. Your job is to analyze the website content and give me up to a maximum of 5 topics or categories that most accurately describe the page, along with a confidence score between 0 and 1 for each topic. 

The confidence score should reflect how certain you are that the topic is relevant to the webpage content, with higher scores (closer to 1) indicating greater confidence.`

	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Company name: %s\n", globalOrg.Name))
	prompt.WriteString(fmt.Sprintf("Company description: %s\n", globalOrg.Description))
	prompt.WriteString(fmt.Sprintf("Webpage url: %s\n", url))
	prompt.WriteString("--- Webpage content --- \n")
	prompt.WriteString(*pageContent)
	promptStr := prompt.String()

	temperature := float32(0.2)
	maxOutputTokens := int32(100)
	topics, err := s.aiService.AskAIForWebpageTopics(ctx, interfaces.AskAIRequest{
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
	if topics == nil {
		span.LogFields(log.String("result", "no topics returned"))
		return nil, nil
	}

	cleanTopics := s.processTopics(topics)

	err = s.postgresRepositories.ScrapedWebpageRepository.SetWebpageTopics(ctx, url, cleanTopics)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	tracing.LogObjectAsJson(span, "result.topics", cleanTopics)
	return cleanTopics, nil
}

func (s *webscraperService) processTopics(topics []data_fields.Topic) []string {
	var response []string
	for _, topic := range topics {
		if topic.Confidence > 0.6 {
			response = append(response, topic.Name)
		}
	}
	return response
}
