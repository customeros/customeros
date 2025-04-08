package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type ForwardEmailReplyCapability struct{}

type ForwardEmailReplyInput struct{}

type ForwardEmailReplyConfig struct{}

type ForwardEmailReplyOutput struct{}

func NewForwardEmailReplyCapability() *ForwardEmailReplyCapability {
	return &ForwardEmailReplyCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[ForwardEmailReplyInput, ForwardEmailReplyOutput, ForwardEmailReplyConfig] = (*ForwardEmailReplyCapability)(nil)
)

func (c *ForwardEmailReplyCapability) Type() enum.AgentCapability {
	return enum.CapabilityForwardEmailReply
}

func (c *ForwardEmailReplyCapability) Name() string {
	return "ForwardEmailReply"
}

func (c *ForwardEmailReplyCapability) NewInput() ForwardEmailReplyInput {
	return ForwardEmailReplyInput{}
}

func (c *ForwardEmailReplyCapability) NewConfig() ForwardEmailReplyConfig {
	return ForwardEmailReplyConfig{}
}

func (c *ForwardEmailReplyCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *ForwardEmailReplyCapability) DefaultActive() bool {
	return true
}

func (c *ForwardEmailReplyCapability) ValidateConfig(config ForwardEmailReplyConfig) error {
	return nil
}

func (c *ForwardEmailReplyCapability) ValidateInput(input ForwardEmailReplyInput) error {
	return nil
}

func (c *ForwardEmailReplyCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ForwardEmailReplyInput, ForwardEmailReplyConfig]) (enum.CapabilityExecutionStatus, ForwardEmailReplyOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ForwardEmailReplyCapability.Execute")
	defer spans.Finish()

	result := ForwardEmailReplyOutput{}

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
