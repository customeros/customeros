package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AgentCapability[I, O, C any] interface {
	AgentCapabilityUntyped
	GetInput() I
	GetConfig() C
	Execute(ctx context.Context, inputData I, configData C) (outputData O, error error)
	ValidateConfig(C) error
	ValidateInput(I) error
}

type AgentCapabilityUntyped interface {
	Type() enum.AgentCapability
	GetConfig() any
}

type AgentCapabilityExecutionService interface {
	Execute(ctx context.Context, capability postgres_entity.Capability, params map[string]any, executors map[enum.AgentCapability]AgentCapabilityUntyped) (map[string]any, error)
}
