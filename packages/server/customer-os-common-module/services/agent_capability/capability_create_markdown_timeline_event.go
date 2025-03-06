package agent_capability

import (
	"context"

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

type CreateMarkdownTimelineEventCapability struct {
	markdownService interfaces.MarkdownEventService
}

type CreateMarkdownTimelineEventInput struct {
	OrganizationID string `json:"organizationId"`
	Message        string `json:"message"`
}

type CreateMarkdownTimelineEventOutput struct {
	MarkdownEventID string `json:"markdownEventId"`
}

func NewCreateMarkdownTimelineEventCapability(markdownService interfaces.MarkdownEventService) *CreateMarkdownTimelineEventCapability {
	return &CreateMarkdownTimelineEventCapability{
		markdownService: markdownService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[CreateMarkdownTimelineEventInput, CreateMarkdownTimelineEventOutput, postgres_entity.NoConfig] = (*CreateMarkdownTimelineEventCapability)(nil)
)

func (c *CreateMarkdownTimelineEventCapability) NewInput() CreateMarkdownTimelineEventInput {
	return CreateMarkdownTimelineEventInput{}
}

func (c *CreateMarkdownTimelineEventCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *CreateMarkdownTimelineEventCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *CreateMarkdownTimelineEventCapability) DefaultActive() bool {
	return true
}

func (c *CreateMarkdownTimelineEventCapability) Type() enum.AgentCapability {
	return enum.CapabilityCreateMarkdownTimelineEvent
}

func (c *CreateMarkdownTimelineEventCapability) Name() string {
	return "Add event to timeline"
}

func (c *CreateMarkdownTimelineEventCapability) ValidateConfig(config postgres_entity.NoConfig) error {
	return nil
}

func (c *CreateMarkdownTimelineEventCapability) ValidateInput(input CreateMarkdownTimelineEventInput) error {
	if input.OrganizationID == "" {
		return errors.New("OrganizationId must be set")
	}
	if input.Message == "" {
		return errors.New("Message must be set")
	}
	return nil
}

func (c *CreateMarkdownTimelineEventCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[CreateMarkdownTimelineEventInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, CreateMarkdownTimelineEventOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateMarkdownTimelineEventCapability.Execute")
	defer span.Finish()
	tracing.SetDefaultAgentCapabilitySpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := CreateMarkdownTimelineEventOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	markdownEventId, err := c.markdownService.Save(ctx, nil, nil, data_fields.MarkdownEventFields{
		OrganizationId: utils.StringPtr(executionContainer.InputData.OrganizationID),
		Content:        utils.StringPtr(executionContainer.InputData.Message),
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	result.MarkdownEventID = markdownEventId
	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
