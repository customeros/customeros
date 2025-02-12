package agent_capability

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type TemplateCapability struct {
	// TODO add services here
}

func NewTemplateCapability() *TemplateCapability {
	return &TemplateCapability{}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[TemplateInput, TemplateOutput, postgres_entity.NoConfig] = (*TemplateCapability)(nil)
)

func (c *TemplateCapability) Type() enum.AgentCapability {
	// TODO change this to the correct capability
	return enum.CapabilityTemplate
}

func (c *TemplateCapability) Name() string {
	return "Template"
}

func (c *TemplateCapability) NewInput() TemplateInput {
	return TemplateInput{}
}

func (c *TemplateCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *TemplateCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *TemplateCapability) ValidateInput(input TemplateInput) error {
	return nil
}

func (c *TemplateCapability) ValidateConfig(config postgres_entity.NoConfig) error {
	return nil
}

type TemplateInput struct {
	// TODO add input here
}

type TemplateOutput struct {
	// TODO add output here
}

func (c *TemplateCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[TemplateInput, postgres_entity.NoConfig]) (bool, TemplateOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TemplateCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := TemplateOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return false, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return false, result, err
	}

	// TODO implement execition here

	tracing.LogObjectAsJson(span, "result", result)
	return true, result, nil
}
