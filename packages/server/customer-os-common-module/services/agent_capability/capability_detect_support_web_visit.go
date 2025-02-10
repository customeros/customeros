package agent_capability

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type DetectSupportWebVisitCapability struct {
	events *events.EventsService
}

type DetectSupportWebVisitInput struct {
	UniquePageViews []string `json:"uniquePageViews"`
	OrganizationId  string   `json:"organizationId"`
}

type DetectSupportWebVisitConfig struct {
	SupportUrls        ConfigMultipleValues `json:"supportUrls"`
	SupportUrlPatterns ConfigMultipleValues `json:"supportUrlPatterns"`
}

func (c *DetectSupportWebVisitConfig) Validate() bool {
	if len(c.SupportUrlPatterns.Value) == 0 && len(c.SupportUrls.Value) == 0 {
		return false
	}

	if len(c.SupportUrls.Value) > 0 {
		cleanUrls := make([]string, 0, len(c.SupportUrls.Value))
		for _, url := range c.SupportUrls.Value {
			cleanUrls = append(cleanUrls, utils.NormalizeUrlPath(url))
		}
		c.SupportUrls.Value = cleanUrls
	}

	if len(c.SupportUrlPatterns.Value) > 0 {
		cleanPatterns := make([]string, 0, len(c.SupportUrlPatterns.Value))
		for _, pattern := range c.SupportUrlPatterns.Value {
			clean := strings.TrimPrefix(pattern, "https://")
			clean = strings.TrimPrefix(clean, "http://")
			clean = strings.TrimPrefix(clean, "www.")
			cleanPatterns = append(cleanPatterns, clean)
		}
		c.SupportUrlPatterns.Value = cleanPatterns
	}

	return true
}

func NewDetectSupportWebVisitCapability(events *events.EventsService) *DetectSupportWebVisitCapability {
	return &DetectSupportWebVisitCapability{
		events: events,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[DetectSupportWebVisitInput, NoOutput, DetectSupportWebVisitConfig] = (*DetectSupportWebVisitCapability)(nil)
)

func (c *DetectSupportWebVisitCapability) Type() enum.AgentCapability {
	return enum.CapabilityDetectSupportWebvisit
}

func (c *DetectSupportWebVisitCapability) Name() string {
	return "Detect support web visit"
}

func (c *DetectSupportWebVisitCapability) NewInput() DetectSupportWebVisitInput {
	return DetectSupportWebVisitInput{}
}

func (c *DetectSupportWebVisitCapability) NewConfig() DetectSupportWebVisitConfig {
	return DetectSupportWebVisitConfig{}
}

func (c *DetectSupportWebVisitCapability) DefaultConfig() any {
	config := c.NewConfig()
	config.SupportUrlPatterns.Value = []string{"**support**"}
	return &config
}

func (c *DetectSupportWebVisitCapability) ValidateConfig(config DetectSupportWebVisitConfig) error {
	if len(config.SupportUrlPatterns.Value) == 0 && len(config.SupportUrls.Value) == 0 {
		return errors.New("Support URL or pattern not configured")
	}
	return nil
}

func (c *DetectSupportWebVisitCapability) ValidateInput(input DetectSupportWebVisitInput) error {
	if len(input.UniquePageViews) == 0 {
		return errors.New("No page views to analyze")
	}
	return nil
}

func (c *DetectSupportWebVisitCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[DetectSupportWebVisitInput, DetectSupportWebVisitConfig]) (bool, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DetectSupportWebVisitCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	result := NoOutput{}

	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return false, result, err
	}
	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return false, result, err
	}

	for _, page := range executionContainer.InputData.UniquePageViews {
		if c.isSupportVisit(page, executionContainer.ConfigData) {
			// throw event
			err := c.events.Publisher.PublishFanoutEvent(ctx, executionContainer.InputData.OrganizationId, model.ORGANIZATION, dto.TagCompany{})
			if err != nil {
				tracing.TraceErr(span, err)
				return true, result, err
			}
			return true, result, nil
		}
	}

	tracing.LogObjectAsJson(span, "result", result)
	return true, result, nil
}

func (c *DetectSupportWebVisitCapability) isSupportVisit(page string, config DetectSupportWebVisitConfig) bool {
	for _, url := range config.SupportUrls.Value {
		if strings.Contains(page, url) {
			return true
		}
		if strings.Contains(url, page) {
			return true
		}
	}

	for _, pattern := range config.SupportUrlPatterns.Value {
		if utils.MatchUrlPattern(pattern, page) {
			return true
		}
	}

	return false
}
