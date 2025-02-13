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

type CreatePaymentLinkCapability struct {
	invoiceService interfaces.InvoiceService
}

func NewCreatePaymentLinkCapability(invoiceService interfaces.InvoiceService) *CreatePaymentLinkCapability {
	return &CreatePaymentLinkCapability{
		invoiceService: invoiceService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[CreatePaymentLinkInput, CreatePaymentLinkOutput, postgres_entity.NoConfig] = (*CreatePaymentLinkCapability)(nil)
)

func (c *CreatePaymentLinkCapability) Type() enum.AgentCapability {
	return enum.CapabilityCreatePaymentLink
}

func (c *CreatePaymentLinkCapability) Name() string {
	return "Create payment link"
}

func (c *CreatePaymentLinkCapability) NewInput() CreatePaymentLinkInput {
	return CreatePaymentLinkInput{}
}

func (c *CreatePaymentLinkCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *CreatePaymentLinkCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *CreatePaymentLinkCapability) ValidateInput(input CreatePaymentLinkInput) error {
	if input.InvoiceID == "" {
		return errors.New("InvoiceID required")
	}
	return nil
}

func (c *CreatePaymentLinkCapability) ValidateConfig(postgres_entity.NoConfig) error {
	return nil
}

type CreatePaymentLinkInput struct {
	InvoiceID string `json:"invoiceId"`
}

type CreatePaymentLinkOutput struct{}

func (c *CreatePaymentLinkCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[CreatePaymentLinkInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, CreatePaymentLinkOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreatePaymentLinkCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := CreatePaymentLinkOutput{}

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

	err = c.invoiceService.GenerateNewPaymentLink(ctx, executionContainer.InputData.InvoiceID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionCompleted, result, err // failed payment link generation is no a blocker
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
