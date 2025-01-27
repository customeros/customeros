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

type SendWebVisitorSlackNotificationCapability struct {
	notificationService interfaces.NotificationService
}

func (c *SendWebVisitorSlackNotificationCapability) ValidateConfig() error {
	//TODO implement me
	panic("implement me")
}

func (c *SendWebVisitorSlackNotificationCapability) ValidateInput() error {
	//TODO implement me
	panic("implement me")
}

func (c *SendWebVisitorSlackNotificationCapability) GetInput() any {
	return &SendSlackNotificationInput{}
}

func (c *SendWebVisitorSlackNotificationCapability) GetConfig() any {
	return &SendSlackNotificationConfig{}
}

func (c *SendWebVisitorSlackNotificationCapability) GetOutput() any {
	return &SendSlackNotificationResult{}
}

func NewSendWebVisitorSlackNotificationCapability(notificationService interfaces.NotificationService) *SendWebVisitorSlackNotificationCapability {
	return &SendWebVisitorSlackNotificationCapability{
		notificationService: notificationService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapabilityExecution[SendSlackNotificationInput, SendSlackNotificationResult, SendSlackNotificationConfig] = (*SendWebVisitorSlackNotificationCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                                                         = (*SendWebVisitorSlackNotificationCapability)(nil)
)

type SendWebVisitorSlackNotificationInput struct {
	Message *string `json:"message,omitempty"`
}

type SendWebVisitorSlackNotificationResult struct {
	Success bool `json:"success"`
}

type SendWebVisitorSlackNotificationConfig struct {
	ChannelID string `json:"channelId"`
}

func (c *SendWebVisitorSlackNotificationCapability) Execute(ctx context.Context, data SendSlackNotificationInput, config SendSlackNotificationConfig) (SendSlackNotificationResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendWebVisitorSlackNotificationCapability.Execute")
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
func (c *SendWebVisitorSlackNotificationCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
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
