package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type AgentCapabilityService interface {
	InitCapabilityRegistry() error
	RegisterCapability(ctx context.Context, capabilityRecord postgres_entity.AgentCapabilityRegistry) error
	ExecuteCapability(ctx context.Context, executionContainer *dto.CapabilityExecutionContainer) error
}
