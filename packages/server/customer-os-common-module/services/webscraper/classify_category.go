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

func (s *webscraperService) ClassifyWebpageCategory(ctx context.Context, url string) (enum.WebpageCategory, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "webscraperService.ClassifyWebpageCategory")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	webpage, err := s.postgresRepositories.ScrapedWebpageRepository.GetWebpage(ctx, url, 365)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
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
	globalOrg, err := s.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, webpage.PrimaryDomain)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("Company name: %s", globalOrg.Name))
	prompt.WriteString(fmt.Sprintf("Company description: %s", globalOrg.Description))
	prompt.WriteString(fmt.Sprintf("Webpage url: %s", webpage.Url))
	if webpage.Content != "" {
		prompt.WriteString("--- Webpage content --- ")
		prompt.WriteString(webpage.Content)
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

	return enum.GetWebpageCategory(*answer), nil
}
