package agent_capability

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type CreateMarkdownTimelineEventCapability struct {
	markdownService interfaces.MarkdownEventService
}

type CreateMarkdownTimelineEventInput struct {
	OrganizationID string `json:"organizationId"`
	Message        string `json:"message"`
}

type CreateMarkdownTimelineEventOutput struct {
	CapabilityOutput
	MarkdownEventID string `json:"markdownEventId"`
}

func NewCreateMarkdownTimelineEventCapability(markdownService interfaces.MarkdownEventService) *CreateMarkdownTimelineEventCapability {
	return &CreateMarkdownTimelineEventCapability{
		markdownService: markdownService,
	}
}

func (c *CreateMarkdownTimelineEventCapability) GetInput() any {
	return &CreateMarkdownTimelineEventInput{}
}

func (c *CreateMarkdownTimelineEventCapability) GetConfig() any {
	return &NoConfig{}
}

func (c *CreateMarkdownTimelineEventCapability) GetOutput() any {
	return &CreateMarkdownTimelineEventOutput{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapabilityExecution[CreateMarkdownTimelineEventInput, CreateMarkdownTimelineEventOutput, NoConfig] = (*CreateMarkdownTimelineEventCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                                                  = (*CreateMarkdownTimelineEventCapability)(nil)
)

func (c *CreateMarkdownTimelineEventCapability) ValidateConfig(config NoConfig) error {
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

func (c *CreateMarkdownTimelineEventCapability) Execute(ctx context.Context, data CreateMarkdownTimelineEventInput, config NoConfig) (CreateMarkdownTimelineEventOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateMarkdownTimelineEventCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", data)
	tracing.LogObjectAsJson(span, "config", config)

	result := CreateMarkdownTimelineEventOutput{}

	if err := c.ValidateInput(data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return result, err
	}
	if err := c.ValidateConfig(config); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return result, err
	}
	result.ExecutionValidated = true

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return result, err
	}

	markdownEventId, err := c.markdownService.Save(ctx, nil, nil, data_fields.MarkdownEventFields{
		OrganizationId: utils.StringPtr(data.OrganizationID),
		Content:        utils.StringPtr(data.Message),
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}

	result.Completed = true
	result.MarkdownEventID = markdownEventId
	tracing.LogObjectAsJson(span, "result", result)
	return result, nil
}

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *CreateMarkdownTimelineEventCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*CreateMarkdownTimelineEventInput)
	if !ok || typedInput == nil {
		return nil, fmt.Errorf("invalid input type: expected CreateMarkdownTimelineEventInput")
	}

	typedConfig, ok := config.(*NoConfig)
	if !ok || typedConfig == nil {
		return nil, fmt.Errorf("invalid config type: expected NoConfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
