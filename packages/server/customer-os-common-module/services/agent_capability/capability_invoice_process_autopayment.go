package agent_capability

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type ProcessAutopaymentCapability struct {
	invoiceService interfaces.InvoiceService
}

func NewProcessAutopaymentCapability(invoiceService interfaces.InvoiceService) *ProcessAutopaymentCapability {
	return &ProcessAutopaymentCapability{
		invoiceService: invoiceService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[ProcessAutopaymentInput, ProcessAutopaymentOutput, ProcessAutopaymentConfig] = (*ProcessAutopaymentCapability)(nil)
)

func (c *ProcessAutopaymentCapability) Type() enum.AgentCapability {
	return enum.CapabilityProcessAutopayment
}

func (c *ProcessAutopaymentCapability) Name() string {
	return "Manage online payment"
}

func (c *ProcessAutopaymentCapability) NewInput() ProcessAutopaymentInput {
	return ProcessAutopaymentInput{}
}

func (c *ProcessAutopaymentCapability) NewConfig() ProcessAutopaymentConfig {
	return ProcessAutopaymentConfig{}
}

func (c *ProcessAutopaymentCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *ProcessAutopaymentCapability) ValidateInput(ProcessAutopaymentInput) error {
	return nil
}

func (c *ProcessAutopaymentCapability) ValidateConfig(ProcessAutopaymentConfig) error {
	return nil
}

type ProcessAutopaymentConfig struct {
	Stripe ConfigSingleBoolValue `json:"stripe"`
}

type ProcessAutopaymentInput struct {
	InvoiceID string `json:"invoiceId"`
}

type ProcessAutopaymentOutput struct{}

func (c *ProcessAutopaymentCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ProcessAutopaymentInput, ProcessAutopaymentConfig]) (enum.CapabilityExecutionStatus, ProcessAutopaymentOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProcessAutopaymentCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := ProcessAutopaymentOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	err := c.invoiceService.AutopayInvoice(ctx, executionContainer.InputData.InvoiceID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
