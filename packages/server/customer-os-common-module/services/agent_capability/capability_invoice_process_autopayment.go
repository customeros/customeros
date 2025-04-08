package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
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

func (c *ProcessAutopaymentCapability) DefaultActive() bool {
	return false
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
	spans, ctx := telemetry.StartServiceSpan(ctx, "ProcessAutopaymentCapability.Execute")
	defer spans.Finish()

	spans.LogObjectAsJson("input", executionContainer.InputData)
	spans.LogObjectAsJson("config", executionContainer.ConfigData)

	result := ProcessAutopaymentOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	err := c.invoiceService.AutopayInvoice(ctx, executionContainer.InputData.InvoiceID)
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionError, result, err
	}

	spans.LogObjectAsJson("result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
