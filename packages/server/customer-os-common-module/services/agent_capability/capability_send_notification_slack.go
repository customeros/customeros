package agent_capability

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
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
var (
	_ interfaces.AgentCapability[SendSlackNotificationInput, NoOutput, SendSlackNotificationConfig] = (*SendSlackNotificationCapability)(nil)
)

func (c *SendSlackNotificationCapability) Type() enum.AgentCapability {
	return enum.CapabilitySendSlackNotification
}

func (c *SendSlackNotificationCapability) Name() string {
	return "Send a slack notification"
}

func (c *SendSlackNotificationCapability) NewInput() SendSlackNotificationInput {
	return SendSlackNotificationInput{}
}

func (c *SendSlackNotificationCapability) NewConfig() SendSlackNotificationConfig {
	return SendSlackNotificationConfig{}
}

func (c *SendSlackNotificationCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *SendSlackNotificationCapability) ValidateConfig(config SendSlackNotificationConfig) error {
	if config.ChannelID.Value == "" {
		return errors.New("ChannelID must be set")
	}
	return nil
}

func (c *SendSlackNotificationCapability) ValidateInput(input SendSlackNotificationInput) error {
	if utils.IfNotNilString(input.Message) == "" {
		return errors.New("Message must be set")
	}
	return nil
}

type SendSlackNotificationInput struct {
	Message string `json:"message"`
}

type SendSlackNotificationConfig struct {
	ChannelID SlackChannelIdConfig `json:"channelId"`
}

type SlackChannelIdConfig struct {
	Value string `json:"value"`
	Error string `json:"error"`
}

func (c *SendSlackNotificationCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[SendSlackNotificationInput, SendSlackNotificationConfig]) (bool, NoOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendSlackNotificationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	result := NoOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return false, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return false, result, err
	}

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set on context")
		tracing.TraceErr(span, err)
		return true, result, err
	}

	err := c.notificationService.NotifySlackChannel(ctx, tenant, executionContainer.ConfigData.ChannelID.Value, &executionContainer.InputData.Message)
	if err != nil {
		tracing.TraceErr(span, err)
		return true, result, err
	}

	tracing.LogObjectAsJson(span, "result", result)
	return true, result, nil
}
