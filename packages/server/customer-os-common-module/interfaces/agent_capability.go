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

type AgentCapabilityExecution[I any, O any] interface {
	Execute(ctx context.Context, inputData I) (outputData O, error error)
}
