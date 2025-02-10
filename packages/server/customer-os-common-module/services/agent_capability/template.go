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

type xxxxxCapability struct{}

type xxxxxInput struct{}

type xxxxxConfig struct{}

type xxxxxOutput struct{}

func NewxxxxxCapability() *xxxxxCapability {
	return &xxxxxCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[xxxxxInput, xxxxxOutput, xxxxxConfig] = (*xxxxxCapability)(nil)
)

func (c *xxxxxCapability) Type() enum.AgentCapability {
	return enum.CapabilityCheckSupportNeed
}

func (c *xxxxxCapability) Name() string {
	return "xxxxx"
}

func (c *xxxxxCapability) NewInput() xxxxxInput {
	return xxxxxInput{}
}

func (c *xxxxxCapability) NewConfig() xxxxxConfig {
	return xxxxxConfig{}
}

func (c *xxxxxCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *xxxxxCapability) ValidateConfig(config xxxxxConfig) error {
	return nil
}

func (c *xxxxxCapability) ValidateInput(input xxxxxInput) error {
	return nil
}

func (c *xxxxxCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[xxxxxInput, xxxxxConfig]) (bool, xxxxxOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "xxxxxCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	result := xxxxxOutput{}

	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return false, result, err
	}
	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return false, result, err
	}

	tracing.LogObjectAsJson(span, "result", result)
	return true, result, nil
}
