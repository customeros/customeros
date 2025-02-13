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

type SendInvoiceViaEmailCapability struct {
	invoiceService interfaces.InvoiceService
}

func NewSendInvoiceViaEmailCapability(invoiceService interfaces.InvoiceService) *SendInvoiceViaEmailCapability {
	return &SendInvoiceViaEmailCapability{
		invoiceService: invoiceService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[SendInvoiceViaEmailInput, SendInvoiceViaEmailOutput, postgres_entity.NoConfig] = (*SendInvoiceViaEmailCapability)(nil)
)

func (c *SendInvoiceViaEmailCapability) Type() enum.AgentCapability {
	return enum.CapabilitySendInvoiceViaEmail
}

func (c *SendInvoiceViaEmailCapability) Name() string {
	return "Send invoice via email"
}

func (c *SendInvoiceViaEmailCapability) NewInput() SendInvoiceViaEmailInput {
	return SendInvoiceViaEmailInput{}
}

func (c *SendInvoiceViaEmailCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *SendInvoiceViaEmailCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *SendInvoiceViaEmailCapability) ValidateInput(input SendInvoiceViaEmailInput) error {
	if input.InvoiceID == "" {
		return errors.New("InvoiceID required")
	}
	return nil
}

func (c *SendInvoiceViaEmailCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

type SendInvoiceViaEmailInput struct {
	InvoiceID string `json:"invoiceId"`
}

type SendInvoiceViaEmailOutput struct{}

func (c *SendInvoiceViaEmailCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[SendInvoiceViaEmailInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, SendInvoiceViaEmailOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendInvoiceViaEmailCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := SendInvoiceViaEmailOutput{}

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

	err = c.invoiceService.SendPayInvoiceNotification(ctx, executionContainer.InputData.InvoiceID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionCompleted, result, err // failed send invoice email is not a blocker, since a new attempt will be made automatically by cron
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
