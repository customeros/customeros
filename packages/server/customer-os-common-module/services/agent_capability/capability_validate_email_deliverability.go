package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type ValidateEmailDeliverabilityCapability struct{}

type ValidateEmailDeliverabilityInput struct{}

type ValidateEmailDeliverabilityConfig struct{}

type ValidateEmailDeliverabilityOutput struct{}

func NewValidateEmailDeliverabilityCapability() *ValidateEmailDeliverabilityCapability {
	return &ValidateEmailDeliverabilityCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[ValidateEmailDeliverabilityInput, ValidateEmailDeliverabilityOutput, ValidateEmailDeliverabilityConfig] = (*ValidateEmailDeliverabilityCapability)(nil)
)

func (c *ValidateEmailDeliverabilityCapability) Type() enum.AgentCapability {
	return enum.CapabilityValidateEmailAddressDeliverability
}

func (c *ValidateEmailDeliverabilityCapability) Name() string {
	return "ValidateEmailDeliverability"
}

func (c *ValidateEmailDeliverabilityCapability) NewInput() ValidateEmailDeliverabilityInput {
	return ValidateEmailDeliverabilityInput{}
}

func (c *ValidateEmailDeliverabilityCapability) NewConfig() ValidateEmailDeliverabilityConfig {
	return ValidateEmailDeliverabilityConfig{}
}

func (c *ValidateEmailDeliverabilityCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *ValidateEmailDeliverabilityCapability) DefaultActive() bool {
	return true
}

func (c *ValidateEmailDeliverabilityCapability) ValidateConfig(config ValidateEmailDeliverabilityConfig) error {
	return nil
}

func (c *ValidateEmailDeliverabilityCapability) ValidateInput(input ValidateEmailDeliverabilityInput) error {
	return nil
}

func (c *ValidateEmailDeliverabilityCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ValidateEmailDeliverabilityInput, ValidateEmailDeliverabilityConfig]) (enum.CapabilityExecutionStatus, ValidateEmailDeliverabilityOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ValidateEmailDeliverabilityCapability.Execute")
	defer spans.Finish()

	spans.TagTenant(common.GetTenantFromContext(ctx))

	result := ValidateEmailDeliverabilityOutput{}

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
