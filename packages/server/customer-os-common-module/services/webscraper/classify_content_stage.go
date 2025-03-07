package webscraper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

func (s *webscraperService) ClassifyContentStage(ctx context.Context, url string, pageContent *string) (enum.CustomerJourneyStage, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "webscraperService.ClassifyContentStage")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	_, primaryDomain := domaincheck.PrimaryDomainCheck(utils.ExtractDomain(url))

	if pageContent == nil {
		webpage, err := s.postgresRepositories.ScrapedWebpageRepository.GetWebpage(ctx, url, 365)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		if webpage == nil {
			err := errors.New("webpage doesn't exist")
			span.LogKV("url", url)
			tracing.TraceErr(span, err)
			return "", err
		}

		pageContent = &webpage.Content
		primaryDomain = webpage.PrimaryDomain
	}

	systemPrompt := `I will provide you with the scraped content of a webpage and a brief description of the company who owns it.  Your job is to analyze the website content and tell me what part of the customer journey the content most closely speaks to.  Your choices are: 
    Problem Recognition - when buyers first identify business challenges or needs and they're researching to understand the problem and it's implications.
    Solution Evaluation - when buyers are actively exploring options to solve their problem.
    Decision Preparation - when buyers are preparing to make a purchase decision and are developing business cases and/or addressing implementation concerns.
    Onboarding - when they're looking at technical setup, user training, and getting started guides.
    Outcome Attainment - when customers are looking at best practices for maximizing early results and reading case studies showcasing similar wins.
    Sustained Success - when they're looking at advanced guides, integrations with other systems, and other business transformation content.`

	globalOrg, err := s.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Company name: %s\n", globalOrg.Name))
	prompt.WriteString(fmt.Sprintf("Company description: %s\n", globalOrg.Description))
	prompt.WriteString(fmt.Sprintf("Webpage url: %s\n", url))
	prompt.WriteString("--- Webpage content --- \n")
	prompt.WriteString(*pageContent)
	promptStr := prompt.String()

	temperature := float32(0.2)
	maxOutputTokens := int32(25)
	stage, err := s.aiService.AskAIForContentStage(ctx, interfaces.AskAIRequest{
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
	if stage == "" {
		return "", nil
	}

	err = s.postgresRepositories.ScrapedWebpageRepository.SetContentStage(ctx, url, stage)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return stage, nil
}
