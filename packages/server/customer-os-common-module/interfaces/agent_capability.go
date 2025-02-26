package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AgentCapability[I, O, C any] interface {
	Execute(ctx context.Context, execution TypedExecutionContainer[I, C]) (status enum.CapabilityExecutionStatus, outputData O, error error)
	ValidateConfig(C) error
	ValidateInput(I) error
	NewInput() I
	NewConfig() C
	AgentCapabilityUntyped
}

type AgentCapabilityUntyped interface {
	Type() enum.AgentCapability
	Name() string
	DefaultConfig() any
	DefaultActive() bool
}

type AgentCapabilityExecutionService interface {
	Execute(ctx context.Context, execution ExecutionContainer) (enum.CapabilityExecutionStatus, map[string]any, error)
}

type ExecutionContainer struct {
	AgentExecutionID string
	Capability       postgres_entity.Capability
	ExecutionParams  map[string]any
	UntypedExecutors map[enum.AgentCapability]AgentCapabilityUntyped
}

type TypedExecutionContainer[I, C any] struct {
	AgentExecutionID string
	InputData        I
	ConfigData       C
}
