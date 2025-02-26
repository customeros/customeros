package agent_capability

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SendPastDueNotificationCapability struct {
	invoiceService interfaces.InvoiceService
}

func NewSendPastDueNotificationCapability(invoiceService interfaces.InvoiceService) *SendPastDueNotificationCapability {
	return &SendPastDueNotificationCapability{
		invoiceService: invoiceService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[SendPastDueNotificationInput, SendPastDueNotificationOutput, postgres_entity.NoConfig] = (*SendPastDueNotificationCapability)(nil)
)

func (c *SendPastDueNotificationCapability) Type() enum.AgentCapability {
	return enum.CapabilitySendPastDueNotification
}

func (c *SendPastDueNotificationCapability) Name() string {
	return "Send a past-due notification"
}

func (c *SendPastDueNotificationCapability) NewInput() SendPastDueNotificationInput {
	return SendPastDueNotificationInput{}
}

func (c *SendPastDueNotificationCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *SendPastDueNotificationCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *SendPastDueNotificationCapability) DefaultActive() bool {
	return true
}

func (c *SendPastDueNotificationCapability) ValidateInput(input SendPastDueNotificationInput) error {
	if input.InvoiceID == "" {
		return errors.New("InvoiceID required")
	}
	return nil
}

func (c *SendPastDueNotificationCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

type SendPastDueNotificationInput struct {
	InvoiceID string `json:"invoiceId"`
}

type SendPastDueNotificationOutput struct{}

func (c *SendPastDueNotificationCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[SendPastDueNotificationInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, SendPastDueNotificationOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendPastDueNotificationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := SendPastDueNotificationOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	invoice, err := c.invoiceService.GetById(ctx, nil, executionContainer.InputData.InvoiceID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}
	if invoice == nil {
		err := errors.New("invoice not found")
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}
	if invoice.DryRun {
		return enum.CapabilityExecutionCompleted, result, nil
	}

	err = c.invoiceService.SendPayReminderInvoiceNotification(ctx, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionCompleted, result, err // failed send invoice email is not a blocker, since a new attempt will be made automatically by cron
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
