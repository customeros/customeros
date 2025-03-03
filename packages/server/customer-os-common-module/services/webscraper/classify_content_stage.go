package webscraper

import (
	"context"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

func (s *webscraperService) ClassifyContentStage(ctx context.Context, url string) (enum.CustomerJourneyStage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "webscraperService.ClassifyContentStage")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	webpage, err := s.postgresRepositories.ScrapedWebpageRepository.GetWebpage(ctx, url, 365)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	systemPrompt := `I will provide you with the scraped content of a webpage and a brief description of the company who owns it.  Your job is to analyze the website content and tell me what part of the customer journey the content most closely speaks to.  Your choices are: 
    Problem Recognition - when buyers first identify business challenges or needs and they're researching to understand the problem and it's implications.
    Solution Evaluation - when buyers are actively exploring options to solve their problem.
    Decision Preparation - when buyers are preparing to make a purchase decision and are developing business cases and/or addressing implementation concerns.
    Onboarding - when they're looking at technical setup, user training, and getting started guides.
    Outcome Attainment - when customers are looking at best practices for maximizing early results and reading case studies showcasing similar wins.
    Sustained Success - when they're looking at advanced guides, integrations with other systems, and other business transformation content.
    Please only respond with the stage name, nothing else.  No commentary or preamble.  Ensure the stage you respond with is either Problem Recognition, Solution Evaluation, Decision Preparation, Onboarding, Outcome Attainment, or Sustained Success.
    `

	globalOrg, err := s.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, webpage.PrimaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Company name: %s", globalOrg.Name))
	prompt.WriteString(fmt.Sprintf("Company description: %s", globalOrg.Description))
	prompt.WriteString(fmt.Sprintf("Webpage url: %s", webpage.Url))
	prompt.WriteString("--- Webpage content --- ")
	prompt.WriteString(webpage.Content)
	promptStr := prompt.String()

	temperature := float32(0.2)
	maxOutputTokens := int32(25)
	answer, err := s.aiService.AskAI(ctx, interfaces.AskAIRequest{
		Model:            enum.AIModelLlama8B,
		SystemPrompt:     &systemPrompt,
		Prompt:           &promptStr,
		ModelTemperature: &temperature,
		MaxOutputTokens:  &maxOutputTokens,
		OutputFormat:     enum.AIOutputText,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if answer == nil {
		return "", nil
	}

	stage, err := s.validateCustomerJourneyStage(ctx, *answer)
	if err != nil {
		stage, err = s.retryClassifyContentStage(ctx, answer)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	err = s.postgresRepositories.ScrapedWebpageRepository.SetContentStage(ctx, url, stage)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return stage, nil
}

func (s *webscraperService) retryClassifyContentStage(ctx context.Context, answer *string) (enum.CustomerJourneyStage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebscraperService.retryClassifyContentStage")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	systemPrompt := `Lets try this again.  I'm going to give you an answer you gave me previously and your job is to clean it up so that you only respond with one of the following strings, nothing else:  Problem Recognition, Solution Evaluation, Decision Preparation, Onboarding, Outcome Attainment, or Sustained Success.`

	temperature := float32(0.1)
	maxOutputTokens := int32(25)
	answer, err := s.aiService.AskAI(ctx, interfaces.AskAIRequest{
		Model:            enum.AIModelLlama8B,
		SystemPrompt:     &systemPrompt,
		Prompt:           answer,
		ModelTemperature: &temperature,
		MaxOutputTokens:  &maxOutputTokens,
		OutputFormat:     enum.AIOutputText,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if answer == nil {
		return "", nil
	}

	return s.validateCustomerJourneyStage(ctx, *answer)
}

func (s *webscraperService) validateCustomerJourneyStage(ctx context.Context, answer string) (enum.CustomerJourneyStage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebscraperService.validateCustomerJourneyStage")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	return enum.GetCustomerJourneyStage(answer)
}
