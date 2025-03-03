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

	systemPrompt := `I will provide you with the url of a webpage and its scraped content (if available), along with a brief description of the company who owns it. Your job is to analyze the website url and content (if available) and return the category that most accurately describe the page. Valid categories are: 
    about
    account
    contact
    help
    legal
    partner
    pricing
    product
    resources
    success story
    other
Please only respond with exactly one of the categories above.  No comments or preamble.  If you are unsure of the category, return other.
`
	globalOrg, err := s.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, primaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Company name: %s", globalOrg.Name))
	prompt.WriteString(fmt.Sprintf("Company description: %s", globalOrg.Description))
	prompt.WriteString(fmt.Sprintf("Webpage url: %s", url))
	if *pageContent != "" {
		prompt.WriteString("--- Webpage content --- ")
		prompt.WriteString(*pageContent)
	}
	promptStr := prompt.String()

	temperature := float32(0.2)
	maxOutputTokens := int32(100)
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

	category := enum.GetWebpageCategory(*answer)
	if category == enum.WebpageUnknown {
		category, err = s.retryClassifyWebpageCategory(ctx, answer)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", nil
		}
	}

	err = s.postgresRepositories.ScrapedWebpageRepository.SetWebpageCategory(ctx, url, category)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return category, nil
}

func (s *webscraperService) retryClassifyWebpageCategory(ctx context.Context, answer *string) (enum.WebpageCategory, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebscraperService.retryClassifyWebpageCategory")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	systemPrompt := `Lets try this again.  I'm going to give you an answer you gave me previously and your job is to clean it up so that you only respond with one of the following strings, nothing else:  about, account, contact, help, legal, partner, pricing, product, resources, success story, or other.`

	temperature := float32(0.1)
	maxOutputTokens := int32(25)
	output, err := s.aiService.AskAI(ctx, interfaces.AskAIRequest{
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
	if output == nil {
		return "", nil
	}

	return enum.GetWebpageCategory(*output), nil
}
