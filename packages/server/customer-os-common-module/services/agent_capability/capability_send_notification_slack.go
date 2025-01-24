package agent_capability

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SendSlackNotificationCapability struct {
	notificationService interfaces.NotificationService
}

func NewSendSlackNotificationCapability(notificationService interfaces.NotificationService) *SendSlackNotificationCapability {
	return &SendSlackNotificationCapability{
		notificationService: notificationService,
	}
}

// Compile-time interface check
var _ interfaces.AgentCapabilityExecution[SendSlackNotificationInput, SendSlackNotificationResult] = (*SendSlackNotificationCapability)(nil)

type SendSlackNotificationInput struct {
	Message   *string
	ChannelID string
}

type SendSlackNotificationResult struct {
	Success bool
}

func (c *SendSlackNotificationCapability) Execute(ctx context.Context, data SendSlackNotificationInput) (SendSlackNotificationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendSlackNotificationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)

	result := SendSlackNotificationResult{}

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		result.Success = false
		return result, err
	}

	err := c.notificationService.NotifySlackChannel(ctx, tenant, data.ChannelID, data.Message)
	if err != nil {
		tracing.TraceErr(span, err)
		result.Success = false
		return result, err
	}

	span.LogKV(
		"event", "capability_executed",
	)
	result.Success = true
	return result, nil
}
