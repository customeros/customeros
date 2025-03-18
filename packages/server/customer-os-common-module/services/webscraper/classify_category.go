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

func (s *webscraperService) ClassifyWebpageCategory(ctx context.Context, url string, pageContent *string) (enum.WebpageCategory, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "webscraperService.ClassifyWebpageCategory")
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

	systemPrompt := `I will provide you with the url of a webpage and its scraped content (if available), along with a brief description of the company who owns it. Your job is to analyze the website url and content (if available) and return the category that most accurately describe the page.`

	globalOrg, err := s.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	var prompt strings.Builder

	if globalOrg != nil {
		prompt.WriteString(fmt.Sprintf("Company name: %s\n", globalOrg.Name))
		prompt.WriteString(fmt.Sprintf("Company description: %s\n", globalOrg.Description))
	} else {
		// If no organization found, just use the domain
		prompt.WriteString(fmt.Sprintf("Website domain: %s\n", primaryDomain))
	}

	prompt.WriteString(fmt.Sprintf("Webpage url: %s\n", url))
	if pageContent != nil && *pageContent != "" {
		prompt.WriteString("--- Webpage content --- \n")
		prompt.WriteString(*pageContent)
	} else {
		prompt.WriteString("--- No content available ---\n")
	}
	promptStr := prompt.String()

	temperature := float32(0.2)
	maxOutputTokens := int32(100)
	category, err := s.aiService.AskAIForWebpageCategory(ctx, interfaces.AskAIRequest{
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

	if category != "" {
		err = s.postgresRepositories.ScrapedWebpageRepository.SetWebpageCategory(ctx, url, category)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	return category, nil
}
