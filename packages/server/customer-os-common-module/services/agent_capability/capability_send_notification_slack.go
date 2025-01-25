package agent_capability

import (
	"context"
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SendSlackNotificationCapability struct {
	notificationService interfaces.NotificationService
}

func (c *SendSlackNotificationCapability) GetInput() any {
	return &SendSlackNotificationInput{}
}

func (c *SendSlackNotificationCapability) GetConfig() any {
	return &SendSlackNotificationConfig{}
}

func (c *SendSlackNotificationCapability) GetOutput() any {
	return &SendSlackNotificationResult{}
}

func NewSendSlackNotificationCapability(notificationService interfaces.NotificationService) *SendSlackNotificationCapability {
	return &SendSlackNotificationCapability{
		notificationService: notificationService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapabilityExecution[SendSlackNotificationInput, SendSlackNotificationResult, SendSlackNotificationConfig] = (*SendSlackNotificationCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                                                         = (*SendSlackNotificationCapability)(nil)
)

type SendSlackNotificationInput struct {
	Message   *string `json:"message,omitempty"`
	ChannelID string  `json:"channel_id"`
}

type SendSlackNotificationResult struct {
	Success bool `json:"success"`
}

type SendSlackNotificationConfig struct {
	ChannelID string `json:"channel_id"`
}

func (c *SendSlackNotificationCapability) Execute(ctx context.Context, data SendSlackNotificationInput, config SendSlackNotificationConfig) (SendSlackNotificationResult, error) {
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

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *SendSlackNotificationCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*SendSlackNotificationInput)
	if !ok {
		return nil, fmt.Errorf("invalid input type: expected SendSlackNotificationInput")
	}

	typedConfig, ok := config.(*SendSlackNotificationConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type: expected SendSlackNotificationConfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
