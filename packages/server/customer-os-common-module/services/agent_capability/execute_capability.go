package agent_capability

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

func (c *agentCapabilityService) ExecuteCapability(ctx context.Context, executionContainer *dto.CapabilityExecutionContainer) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.ExecuteCapability")
	defer span.Finish()
	tracing.TagComponentService(span)

	span.LogKV(
		"event", "execute_capability",
		"agent_id", executionContainer.AgentID,
		"capability", executionContainer.Capability.String(),
		"input_type", fmt.Sprintf("%T", executionContainer.InputData),
	)

	if executionContainer.AgentID == "" {
		err := errors.New("AgentID not set")
		tracing.TraceErr(span, err)
		return err
	}

	handler, exists := c.executionHandlers[executionContainer.Capability]
	if !exists {
		err := fmt.Errorf("capability %s not configured for execution", executionContainer.Capability)
		tracing.TraceErr(span, err)
		return err
	}

	return handler(ctx, executionContainer)
}
