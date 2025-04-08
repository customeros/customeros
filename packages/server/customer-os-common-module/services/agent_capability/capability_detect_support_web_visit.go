package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"strings"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

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
	SupportUrls ConfigMultipleValues `json:"supportUrls"`
}

func (c *DetectSupportWebVisitConfig) Validate() bool {
	if len(c.SupportUrls.Value) == 0 {
		return false
	}

	if len(c.SupportUrls.Value) > 0 {
		cleanUrls := make([]string, 0, len(c.SupportUrls.Value))
		for _, url := range c.SupportUrls.Value {
			cleanUrls = append(cleanUrls, utils.NormalizeUrlPath(url))
		}
		c.SupportUrls.Value = cleanUrls
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
	config.SupportUrls.Value = []string{}
	return &config
}

func (c *DetectSupportWebVisitCapability) DefaultActive() bool {
	return true
}

func (c *DetectSupportWebVisitCapability) ValidateConfig(config DetectSupportWebVisitConfig) error {
	if len(config.SupportUrls.Value) == 0 {
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

func (c *DetectSupportWebVisitCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[DetectSupportWebVisitInput, DetectSupportWebVisitConfig]) (enum.CapabilityExecutionStatus, NoOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "DetectSupportWebVisitCapability.Execute")
	defer spans.Finish()

	result := NoOutput{}

	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}

	for _, page := range executionContainer.InputData.UniquePageViews {
		if c.isSupportVisit(page, executionContainer.ConfigData) {
			// throw event
			err := c.events.Publisher.PublishFanoutEvent(ctx, executionContainer.InputData.OrganizationId, model.ORGANIZATION, dto.CompanyNeedsHelp{})
			if err != nil {
				spans.TraceError(err)
				return enum.CapabilityExecutionError, result, err
			}
			return enum.CapabilityExecutionError, result, nil
		}
	}

	spans.LogObjectAsJson("result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}

func (c *DetectSupportWebVisitCapability) isSupportVisit(page string, config DetectSupportWebVisitConfig) bool {
	for _, url := range config.SupportUrls.Value {
		if strings.Contains(page, url) {
			return true
		}
	}
	return false
}
