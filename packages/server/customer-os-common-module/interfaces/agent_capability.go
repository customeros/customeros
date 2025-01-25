package interfaces

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type AgentCapabilityRegistry interface {
	RegisterCapability(ctx context.Context, capabilityRecord postgres_entity.AgentCapabilityRegistry) error
	ExecuteCapability(ctx context.Context, executionContainer *dto.CapabilityExecutionContainer) error
}

type AgentCapabilityExecution[I any, O any, C any] interface {
	Execute(ctx context.Context, inputData I, configData C) (outputData O, error error)
}

type AgentCapabilityUntyped interface {
	Capability
	ExecuteUntyped(ctx context.Context, input any, config any) (any, error)
}

type Capability interface {
	GetInput() any
	GetConfig() any
	GetOutput() any
}
