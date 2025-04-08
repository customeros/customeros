package agent_capability

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
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

func (c *TEMPLATECapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[TEMPLATEInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, TEMPLATEOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TEMPLATECapability.Execute")
	defer spans.Finish()

	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	spans.LogObjectAsJson("input", executionContainer.InputData)
	spans.LogObjectAsJson("config", executionContainer.ConfigData)

	result := TEMPLATEOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid input"))
		return false, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid config"))
		return false, result, err
	}

	// TODO implement execution here

	spans.LogObjectAsJson("result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
