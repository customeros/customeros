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

type EnrichEmailAddressCapability struct{}

type EnrichEmailAddressInput struct{}

type EnrichEmailAddressConfig struct{}

type EnrichEmailAddressOutput struct{}

func NewEnrichEmailAddressCapability() *EnrichEmailAddressCapability {
	return &EnrichEmailAddressCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[EnrichEmailAddressInput, EnrichEmailAddressOutput, EnrichEmailAddressConfig] = (*EnrichEmailAddressCapability)(nil)
)

func (c *EnrichEmailAddressCapability) Type() enum.AgentCapability {
	return enum.CapabilityEnrichEmailAddress
}

func (c *EnrichEmailAddressCapability) Name() string {
	return "EnrichEmailAddress"
}

func (c *EnrichEmailAddressCapability) NewInput() EnrichEmailAddressInput {
	return EnrichEmailAddressInput{}
}

func (c *EnrichEmailAddressCapability) NewConfig() EnrichEmailAddressConfig {
	return EnrichEmailAddressConfig{}
}

func (c *EnrichEmailAddressCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *EnrichEmailAddressCapability) ValidateConfig(config EnrichEmailAddressConfig) error {
	return nil
}

func (c *EnrichEmailAddressCapability) ValidateInput(input EnrichEmailAddressInput) error {
	return nil
}

func (c *EnrichEmailAddressCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[EnrichEmailAddressInput, EnrichEmailAddressConfig]) (bool, EnrichEmailAddressOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EnrichEmailAddressCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))

	result := EnrichEmailAddressOutput{}

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
