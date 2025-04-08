package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type SummarizeMessageInput struct {
	EntityId   string           `json:"entityId"`
	EntityType model.EntityType `json:"entityType"`
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

func (c *SummarizeMessageCapability) DefaultActive() bool {
	return true
}

func (c *SummarizeMessageCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

func (c *SummarizeMessageCapability) ValidateInput(input SummarizeMessageInput) error {
	if input.EntityId == "" {
		return errors.New("entityId is required")
	}
	if input.EntityType == "" {
		return errors.New("entityType is required")
	}
	return nil
}

func (c *SummarizeMessageCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[SummarizeMessageInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, NoOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "SummarizeMessageCapability.Execute")
	defer spans.Finish()

	spans.TagTenant(common.GetTenantFromContext(ctx))
	spans.LogObjectAsJson("executionContainer", executionContainer)

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, NoOutput{}, err
	}

	//todo implement

	return enum.CapabilityExecutionCompleted, NoOutput{}, nil
}
