package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SummarizeThreadInput struct {
	EntityId   string           `json:"entityId"`
	EntityType model.EntityType `json:"entityType"`
}

type SummarizeThreadCapability struct {
}

func NewSummarizeThreadCapability() *SummarizeThreadCapability {
	return &SummarizeThreadCapability{}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[SummarizeThreadInput, NoOutput, postgres_entity.NoConfig] = (*SummarizeThreadCapability)(nil)
)

func (c *SummarizeThreadCapability) Type() enum.AgentCapability {
	return enum.CapabilitySummarizeThread
}

func (c *SummarizeThreadCapability) Name() string {
	return "Summarize thread"
}

func (c *SummarizeThreadCapability) NewInput() SummarizeThreadInput {
	return SummarizeThreadInput{}
}

func (c *SummarizeThreadCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *SummarizeThreadCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *SummarizeThreadCapability) DefaultActive() bool {
	return true
}

func (c *SummarizeThreadCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

func (c *SummarizeThreadCapability) ValidateInput(input SummarizeThreadInput) error {
	if input.EntityId == "" {
		return errors.New("entityId is required")
	}
	if input.EntityType == "" {
		return errors.New("entityType is required")
	}
	return nil
}

func (c *SummarizeThreadCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[SummarizeThreadInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SummarizeThreadCapability.Execute")
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
