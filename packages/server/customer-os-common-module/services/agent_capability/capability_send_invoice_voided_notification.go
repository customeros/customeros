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
	return "Send invoice via email"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendInvoiceVoidedNotificationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := SendInvoiceVoidedNotificationOutput{}

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

	err = c.invoiceService.SendVoidedInvoiceNotification(ctx, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionCompleted, result, err // failed send invoice email is not a blocker, since a new attempt will be made automatically by cron
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
