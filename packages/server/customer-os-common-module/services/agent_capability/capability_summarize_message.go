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

type SummarizeMessageInput struct {
}

type SummarizeMessageCapability struct {
}

func NewSummarizeMessageCapability() *SummarizeMessageCapability {
	return &SummarizeMessageCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[SummarizeMessageInput, NoOutput, postgres_entity.NoConfig] = (*SummarizeMessageCapability)(nil)
)

func (c *SummarizeMessageCapability) Type() enum.AgentCapability {
	return enum.CapabilitySummarizeMessage
}

func (c *SummarizeMessageCapability) Name() string {
	return "Summarize message"
}

func (c *SummarizeMessageCapability) NewInput() SummarizeMessageInput {
	return SummarizeMessageInput{}
}

func (c *SummarizeMessageCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *SummarizeMessageCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *SummarizeMessageCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

func (c *SummarizeMessageCapability) ValidateInput(input SummarizeMessageInput) error {
	//todo
	return nil
}

func (c *SummarizeMessageCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[SummarizeMessageInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SummarizeMessageCapability.Execute")
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
