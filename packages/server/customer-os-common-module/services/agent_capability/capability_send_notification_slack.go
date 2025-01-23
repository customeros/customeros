package agent_capability

import (
	"context"
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SendSlackNotificationInput struct {
	Message   *string
	ChannelID string
}

type SendSlackNotificationResult struct{}

func (c *agentCapabilityService) handleSendSlackNotificationExecution(ctx context.Context, executionContainer *dto.CapabilityExecutionContainer) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentCapabilityService.handleSendSlackNotificationExecution")
	defer span.Finish()
	tracing.TagComponentService(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return err
	}

	input, ok := executionContainer.InputData.(SendSlackNotificationInput)
	if !ok {
		err := fmt.Errorf("expected SendSlackNotificationInput, got %T", executionContainer.InputData)
		tracing.TraceErr(span, err)
		return err
	}
	err := c.notificationService.NotifySlackChannel(ctx, tenant, input.ChannelID, input.Message)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	span.LogKV(
		"event", "capability_executed",
	)
	return nil
}
