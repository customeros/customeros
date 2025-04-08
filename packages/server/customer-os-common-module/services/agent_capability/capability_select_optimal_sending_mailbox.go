package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type SelectOptimalSendingMailboxCapability struct{}

type SelectOptimalSendingMailboxInput struct{}

type SelectOptimalSendingMailboxConfig struct{}

type SelectOptimalSendingMailboxOutput struct{}

func NewSelectOptimalSendingMailboxCapability() *SelectOptimalSendingMailboxCapability {
	return &SelectOptimalSendingMailboxCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[SelectOptimalSendingMailboxInput, SelectOptimalSendingMailboxOutput, SelectOptimalSendingMailboxConfig] = (*SelectOptimalSendingMailboxCapability)(nil)
)

func (c *SelectOptimalSendingMailboxCapability) Type() enum.AgentCapability {
	return enum.CapabilitySelectOptimalSendingMailbox
}

func (c *SelectOptimalSendingMailboxCapability) Name() string {
	return "SelectOptimalSendingMailbox"
}

func (c *SelectOptimalSendingMailboxCapability) NewInput() SelectOptimalSendingMailboxInput {
	return SelectOptimalSendingMailboxInput{}
}

func (c *SelectOptimalSendingMailboxCapability) NewConfig() SelectOptimalSendingMailboxConfig {
	return SelectOptimalSendingMailboxConfig{}
}

func (c *SelectOptimalSendingMailboxCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *SelectOptimalSendingMailboxCapability) DefaultActive() bool {
	return true
}

func (c *SelectOptimalSendingMailboxCapability) ValidateConfig(config SelectOptimalSendingMailboxConfig) error {
	return nil
}

func (c *SelectOptimalSendingMailboxCapability) ValidateInput(input SelectOptimalSendingMailboxInput) error {
	return nil
}

func (c *SelectOptimalSendingMailboxCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[SelectOptimalSendingMailboxInput, SelectOptimalSendingMailboxConfig]) (enum.CapabilityExecutionStatus, SelectOptimalSendingMailboxOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "SelectOptimalSendingMailboxCapability.Execute")
	defer spans.Finish()

	result := SelectOptimalSendingMailboxOutput{}

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
