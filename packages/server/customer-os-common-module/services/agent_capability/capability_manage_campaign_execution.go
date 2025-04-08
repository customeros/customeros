package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
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

func (c *ManageCampaignExecutionCapability) DefaultActive() bool {
	return true
}

func (c *ManageCampaignExecutionCapability) ValidateConfig(config ManageCampaignExecutionConfig) error {
	return nil
}

func (c *ManageCampaignExecutionCapability) ValidateInput(input ManageCampaignExecutionInput) error {
	return nil
}

func (c *ManageCampaignExecutionCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ManageCampaignExecutionInput, ManageCampaignExecutionConfig]) (enum.CapabilityExecutionStatus, ManageCampaignExecutionOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ManageCampaignExecutionCapability.Execute")
	defer spans.Finish()

	result := ManageCampaignExecutionOutput{}

	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}

	spans.LogObjectAsJson("result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
