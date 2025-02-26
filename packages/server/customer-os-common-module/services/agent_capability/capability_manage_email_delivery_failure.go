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

type ManageEmailDeliveryFailureCapability struct{}

type ManageEmailDeliveryFailureInput struct{}

type ManageEmailDeliveryFailureConfig struct{}

type ManageEmailDeliveryFailureOutput struct{}

func NewManageEmailDeliveryFailureCapability() *ManageEmailDeliveryFailureCapability {
	return &ManageEmailDeliveryFailureCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[ManageEmailDeliveryFailureInput, ManageEmailDeliveryFailureOutput, ManageEmailDeliveryFailureConfig] = (*ManageEmailDeliveryFailureCapability)(nil)
)

func (c *ManageEmailDeliveryFailureCapability) Type() enum.AgentCapability {
	return enum.CapabilityManageEmailDeliveryFailure
}

func (c *ManageEmailDeliveryFailureCapability) Name() string {
	return "ManageEmailDeliveryFailure"
}

func (c *ManageEmailDeliveryFailureCapability) NewInput() ManageEmailDeliveryFailureInput {
	return ManageEmailDeliveryFailureInput{}
}

func (c *ManageEmailDeliveryFailureCapability) NewConfig() ManageEmailDeliveryFailureConfig {
	return ManageEmailDeliveryFailureConfig{}
}

func (c *ManageEmailDeliveryFailureCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *ManageEmailDeliveryFailureCapability) DefaultActive() bool {
	return true
}

func (c *ManageEmailDeliveryFailureCapability) ValidateConfig(config ManageEmailDeliveryFailureConfig) error {
	return nil
}

func (c *ManageEmailDeliveryFailureCapability) ValidateInput(input ManageEmailDeliveryFailureInput) error {
	return nil
}

func (c *ManageEmailDeliveryFailureCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ManageEmailDeliveryFailureInput, ManageEmailDeliveryFailureConfig]) (enum.CapabilityExecutionStatus, ManageEmailDeliveryFailureOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ManageEmailDeliveryFailureCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	result := ManageEmailDeliveryFailureOutput{}

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
