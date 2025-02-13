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

type ClassifyEmailInput struct {
	RawEmailId string `json:"rawEmailId"`
}

type ClassifyEmailCapability struct {
}

func NewClassifyEmailCapability() *ClassifyEmailCapability {
	return &ClassifyEmailCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[ClassifyEmailInput, NoOutput, postgres_entity.NoConfig] = (*ClassifyEmailCapability)(nil)
)

func (c *ClassifyEmailCapability) Type() enum.AgentCapability {
	return enum.CapabilityClassifyEmail
}

func (c *ClassifyEmailCapability) Name() string {
	return "Classify email"
}

func (c *ClassifyEmailCapability) NewInput() ClassifyEmailInput {
	return ClassifyEmailInput{}
}

func (c *ClassifyEmailCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *ClassifyEmailCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *ClassifyEmailCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

func (c *ClassifyEmailCapability) ValidateInput(input ClassifyEmailInput) error {
	//todo
	return nil
}

func (c *ClassifyEmailCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ClassifyEmailInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ClassifyEmailCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}

	//todo implement

	return enum.CapabilityExecutionCompleted, NoOutput{}, nil
}
