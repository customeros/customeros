package agent_capability

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type ManageCampaignExecutionCapability struct{}

type ManageCampaignExecutionInput struct{}

type ManageCampaignExecutionConfig struct{}

type ManageCampaignExecutionOutput struct{}

func NewManageCampaignExecutionCapability() *ManageCampaignExecutionCapability {
	return &ManageCampaignExecutionCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[ManageCampaignExecutionInput, ManageCampaignExecutionOutput, ManageCampaignExecutionConfig] = (*ManageCampaignExecutionCapability)(nil)
)

func (c *ManageCampaignExecutionCapability) Type() enum.AgentCapability {
	return enum.CapabilityManageCampaignExecution
}

func (c *ManageCampaignExecutionCapability) Name() string {
	return "ManageCampaignExecution"
}

func (c *ManageCampaignExecutionCapability) NewInput() ManageCampaignExecutionInput {
	return ManageCampaignExecutionInput{}
}

func (c *ManageCampaignExecutionCapability) NewConfig() ManageCampaignExecutionConfig {
	return ManageCampaignExecutionConfig{}
}

func (c *ManageCampaignExecutionCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *ManageCampaignExecutionCapability) ValidateConfig(config ManageCampaignExecutionConfig) error {
	return nil
}

func (c *ManageCampaignExecutionCapability) ValidateInput(input ManageCampaignExecutionInput) error {
	return nil
}

func (c *ManageCampaignExecutionCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ManageCampaignExecutionInput, ManageCampaignExecutionConfig]) (enum.CapabilityExecutionStatus, ManageCampaignExecutionOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ManageCampaignExecutionCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	result := ManageCampaignExecutionOutput{}

	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
