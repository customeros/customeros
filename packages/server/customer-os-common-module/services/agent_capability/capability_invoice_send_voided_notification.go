package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type SendInvoiceVoidedNotificationCapability struct {
	invoiceService interfaces.InvoiceService
}

func NewSendInvoiceVoidedNotificationCapability(invoiceService interfaces.InvoiceService) *SendInvoiceVoidedNotificationCapability {
	return &SendInvoiceVoidedNotificationCapability{
		invoiceService: invoiceService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[SendInvoiceVoidedNotificationInput, SendInvoiceVoidedNotificationOutput, postgres_entity.NoConfig] = (*SendInvoiceVoidedNotificationCapability)(nil)
)

func (c *SendInvoiceVoidedNotificationCapability) Type() enum.AgentCapability {
	return enum.CapabilitySendInvoiceVoidedNotification
}

func (c *SendInvoiceVoidedNotificationCapability) Name() string {
	return "Send invoice voided notification"
}

func (c *SendInvoiceVoidedNotificationCapability) NewInput() SendInvoiceVoidedNotificationInput {
	return SendInvoiceVoidedNotificationInput{}
}

func (c *SendInvoiceVoidedNotificationCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *SendInvoiceVoidedNotificationCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *SendInvoiceVoidedNotificationCapability) DefaultActive() bool {
	return true
}

func (c *SendInvoiceVoidedNotificationCapability) ValidateInput(input SendInvoiceVoidedNotificationInput) error {
	if input.InvoiceID == "" {
		return errors.New("InvoiceID required")
	}
	return nil
}

func (c *SendInvoiceVoidedNotificationCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

type SendInvoiceVoidedNotificationInput struct {
	InvoiceID string `json:"invoiceId"`
}

type SendInvoiceVoidedNotificationOutput struct{}

func (c *SendInvoiceVoidedNotificationCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[SendInvoiceVoidedNotificationInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, SendInvoiceVoidedNotificationOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "SendInvoiceVoidedNotificationCapability.Execute")
	defer spans.Finish()

	spans.LogObjectAsJson("input", executionContainer.InputData)
	spans.LogObjectAsJson("config", executionContainer.ConfigData)

	result := SendInvoiceVoidedNotificationOutput{}

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

	err = c.invoiceService.SendVoidedInvoiceNotification(ctx, invoice.Id)
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionCompleted, result, err // failed send invoice email is not a blocker, since a new attempt will be made automatically by cron
	}

	spans.LogObjectAsJson("result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
