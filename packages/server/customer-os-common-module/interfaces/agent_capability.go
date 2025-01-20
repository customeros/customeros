package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AgentCapabilityService interface {
	RegisterCapability(ctx context.Context, agentID string, capabilityType string) error
	DeregisterCapability(ctx context.Context, agentID string, capabilityType string) error
	ExecuteCapability(ctx context.Context, agentID string, capability enum.AgentCapabilityType, input any) (any, error)
}
