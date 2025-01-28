package agent_capability

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SendSlackNotificationCapability struct {
	notificationService interfaces.NotificationService
}

func (c *SendSlackNotificationCapability) ValidateConfig(config SendSlackNotificationConfig) error {
	if config.ChannelID.Value == "" {
		return errors.New("ChannelID must be set")
	}
	return nil
}

func (c *SendSlackNotificationCapability) ValidateInput(input SendSlackNotificationInput) error {
	if input.Message == nil || utils.IfNotNilString(input.Message) == "" {
		return errors.New("Message must be set")
	}
	return nil
}

func (c *SendSlackNotificationCapability) GetInput() any {
	return &SendSlackNotificationInput{}
}

func (c *SendSlackNotificationCapability) GetConfig() any {
	return &SendSlackNotificationConfig{}
}

func (c *SendSlackNotificationCapability) GetOutput() any {
	return &SendSlackNotificationOutput{}
}

func NewSendSlackNotificationCapability(notificationService interfaces.NotificationService) *SendSlackNotificationCapability {
	return &SendSlackNotificationCapability{
		notificationService: notificationService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapabilityExecution[SendSlackNotificationInput, SendSlackNotificationOutput, SendSlackNotificationConfig] = (*SendSlackNotificationCapability)(nil)
	_ interfaces.AgentCapabilityUntyped                                                                                         = (*SendSlackNotificationCapability)(nil)
)

type SendSlackNotificationInput struct {
	Message *string `json:"message,omitempty"`
}

type SendSlackNotificationOutput struct {
	CapabilityOutput
}

type SendSlackNotificationConfig struct {
	ChannelID SlackChannelIdConfig `json:"channelId"`
}

type SlackChannelIdConfig struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

func (c *SendSlackNotificationCapability) Execute(ctx context.Context, data SendSlackNotificationInput, config SendSlackNotificationConfig) (SendSlackNotificationOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendSlackNotificationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", data)
	tracing.LogObjectAsJson(span, "config", config)

	result := SendSlackNotificationOutput{}

	if err := c.ValidateInput(data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return result, err
	}
	if err := c.ValidateConfig(config); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return result, err
	}
	result.ExecutionValidated = true

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return result, err
	}

	err := c.notificationService.NotifySlackChannel(ctx, tenant, config.ChannelID.Value, data.Message)
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}

	result.Completed = true
	tracing.LogObjectAsJson(span, "result", result)
	return result, nil
}

// ExecuteUntyped implements the AgentCapabilityUntyped interface.
// It casts the generic input and config to the specific types and delegates to the typed Execute method.
func (c *SendSlackNotificationCapability) ExecuteUntyped(ctx context.Context, input any, config any) (any, error) {
	typedInput, ok := input.(*SendSlackNotificationInput)
	if !ok || typedInput == nil {
		return nil, fmt.Errorf("invalid input type: expected SendSlackNotificationInput")
	}

	typedConfig, ok := config.(*SendSlackNotificationConfig)
	if !ok || typedConfig == nil {
		return nil, fmt.Errorf("invalid config type: expected SendSlackNotificationConfig")
	}

	return c.Execute(ctx, *typedInput, *typedConfig)
}
