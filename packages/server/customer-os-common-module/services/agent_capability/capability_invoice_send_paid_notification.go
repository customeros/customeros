package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type SendPaidNotificationCapability struct {
	invoiceService interfaces.InvoiceService
}

func NewSendPaidNotificationCapability(invoiceService interfaces.InvoiceService) *SendPaidNotificationCapability {
	return &SendPaidNotificationCapability{
		invoiceService: invoiceService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[SendPaidNotificationInput, SendPaidNotificationOutput, postgres_entity.NoConfig] = (*SendPaidNotificationCapability)(nil)
)

func (c *SendPaidNotificationCapability) Type() enum.AgentCapability {
	return enum.CapabilitySendPaidNotification
}

func (c *SendPaidNotificationCapability) Name() string {
	return "Send a paid notification"
}

func (c *SendPaidNotificationCapability) NewInput() SendPaidNotificationInput {
	return SendPaidNotificationInput{}
}

func (c *SendPaidNotificationCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *SendPaidNotificationCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *SendPaidNotificationCapability) DefaultActive() bool {
	return true
}

func (c *SendPaidNotificationCapability) ValidateInput(input SendPaidNotificationInput) error {
	if input.InvoiceID == "" {
		return errors.New("InvoiceID required")
	}
	return nil
}

func (c *SendPaidNotificationCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

type SendPaidNotificationInput struct {
	InvoiceID string `json:"invoiceId"`
}

type SendPaidNotificationOutput struct{}

func (c *SendPaidNotificationCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[SendPaidNotificationInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, SendPaidNotificationOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "SendPaidNotificationCapability.Execute")
	defer spans.Finish()

	spans.LogObjectAsJson("input", executionContainer.InputData)
	spans.LogObjectAsJson("config", executionContainer.ConfigData)

	result := SendPaidNotificationOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	invoice, err := c.invoiceService.GetById(ctx, nil, executionContainer.InputData.InvoiceID)
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionError, result, err
	}
	if invoice == nil {
		err := errors.New("invoice not found")
		spans.TraceError(err)
		return enum.CapabilityExecutionError, result, err
	}
	if invoice.DryRun {
		return enum.CapabilityExecutionCompleted, result, nil
	}

	err = c.invoiceService.SendPaidInvoiceNotification(ctx, invoice.Id)
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionCompleted, result, err // failed send invoice email is not a blocker, since a new attempt will be made automatically by cron
	}

	spans.LogObjectAsJson("result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
