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

func (c *ForwardEmailReplyCapability) ValidateConfig(config ForwardEmailReplyConfig) error {
	return nil
}

func (c *ForwardEmailReplyCapability) ValidateInput(input ForwardEmailReplyInput) error {
	return nil
}

func (c *ForwardEmailReplyCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ForwardEmailReplyInput, ForwardEmailReplyConfig]) (enum.CapabilityExecutionStatus, ForwardEmailReplyOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ForwardEmailReplyCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	result := ForwardEmailReplyOutput{}

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
