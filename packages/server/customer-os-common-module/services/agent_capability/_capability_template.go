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

type TEMPLATECapability struct {
	// TODO add services here
}

func NewTEMPLATECapability() *TEMPLATECapability {
	return &TEMPLATECapability{}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[TEMPLATEInput, TEMPLATEOutput, postgres_entity.NoConfig] = (*TEMPLATECapability)(nil)
)

func (c *TEMPLATECapability) Type() enum.AgentCapability {
	// TODO change this to the correct capability
	return enum.CapabilityTEMPLATE
}

func (c *TEMPLATECapability) Name() string {
	return "TEMPLATE"
}

func (c *TEMPLATECapability) NewInput() TEMPLATEInput {
	return TEMPLATEInput{}
}

func (c *TEMPLATECapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *TEMPLATECapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *TEMPLATECapability) ValidateInput(input TEMPLATEInput) error {
	return nil
}

func (c *TEMPLATECapability) ValidateConfig(config postgres_entity.NoConfig) error {
	return nil
}

type TEMPLATEInput struct {
	// TODO add input here
}

type TEMPLATEOutput struct {
	// TODO add output here
}

func (c *TEMPLATECapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[TEMPLATEInput, postgres_entity.NoConfig]) (bool, TEMPLATEOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TEMPLATECapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := TEMPLATEOutput{}

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
