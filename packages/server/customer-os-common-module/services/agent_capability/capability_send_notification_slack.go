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

func (c *SendSlackNotificationCapability) Type() enum.AgentCapability {
	return enum.CapabilitySendSlackNotification
}

func (c *SendSlackNotificationCapability) GetInput() SendSlackNotificationInput {
	return SendSlackNotificationInput{}
}

func (c *SendSlackNotificationCapability) GetConfig() SendSlackNotificationConfig {
	return SendSlackNotificationConfig{}
}

func (c *SendSlackNotificationCapability) GetOutput() SendSlackNotificationOutput {
	return SendSlackNotificationOutput{}
}

func NewSendSlackNotificationCapability(notificationService interfaces.NotificationService) *SendSlackNotificationCapability {
	return &SendSlackNotificationCapability{
		notificationService: notificationService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[SendSlackNotificationInput, SendSlackNotificationOutput, SendSlackNotificationConfig] = (*SendSlackNotificationCapability)(nil)
)

type SendSlackNotificationInput struct {
	Message string `json:"message"`
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

	err := c.notificationService.NotifySlackChannel(ctx, tenant, config.ChannelID.Value, &data.Message)
	if err != nil {
		tracing.TraceErr(span, err)
		return result, err
	}

	result.Completed = true
	tracing.LogObjectAsJson(span, "result", result)
	return result, nil
}
