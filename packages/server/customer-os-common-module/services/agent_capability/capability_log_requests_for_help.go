package agent_capability

import (
	"context"
	"fmt"
	"strings"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type LogRequestsForHelpCapability struct {
	markdownEventService interfaces.MarkdownEventService
}

func NewLogRequestsForHelpCapability(markdownEventService interfaces.MarkdownEventService) *LogRequestsForHelpCapability {
	return &LogRequestsForHelpCapability{
		markdownEventService: markdownEventService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[LogRequestsForHelpInput, NoOutput, postgres_entity.NoConfig] = (*LogRequestsForHelpCapability)(nil)
)

func (c *LogRequestsForHelpCapability) Type() enum.AgentCapability {
	return enum.CapabilityLogRequestsForHelp
}

func (c *LogRequestsForHelpCapability) Name() string {
	return "Log requests for help on company timeline"
}

func (c *LogRequestsForHelpCapability) NewInput() LogRequestsForHelpInput {
	return LogRequestsForHelpInput{}
}

func (c *LogRequestsForHelpCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *LogRequestsForHelpCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *LogRequestsForHelpCapability) DefaultActive() bool {
	return true
}

func (c *LogRequestsForHelpCapability) ValidateInput(input LogRequestsForHelpInput) error {
	if input.OrganizationID == "" {
		return errors.New("OrganizationID required")
	}
	if len(input.HelpNeeded) == 0 {
		return errors.New("No values set for help needed")
	}
	return nil
}

func (c *LogRequestsForHelpCapability) ValidateConfig(config postgres_entity.NoConfig) error {
	return nil
}

type LogRequestsForHelpInput struct {
	OrganizationID string   `json:"organizationId"`
	HelpNeeded     []string `json:"helpNeeded"`
}

func (c *LogRequestsForHelpCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[LogRequestsForHelpInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogRequestsForHelpCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := NoOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	timelineEvent := c.createMarkdownTimelineEvent(ctx, executionContainer.InputData)
	appSource := "Agent"
	_, err := c.markdownEventService.Save(ctx, nil, nil, data_fields.MarkdownEventFields{
		AppSource:      &appSource,
		CreatedAt:      utils.NowPtr(),
		OrganizationId: &executionContainer.InputData.OrganizationID,
		Content:        &timelineEvent,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}

func (c *LogRequestsForHelpCapability) createMarkdownTimelineEvent(ctx context.Context, input LogRequestsForHelpInput) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogRequestsForHelpCapability.createMarkdownTimelineEvent")
	defer span.Finish()
	tracing.TagComponentService(span)

	var b strings.Builder

	// Header
	b.WriteString("# Help Needed\n")

	// Help Needed
	for _, p := range input.HelpNeeded {
		b.WriteString(fmt.Sprintf("* %s\n", p))
	}
	b.WriteString("\n")

	return b.String()
}
