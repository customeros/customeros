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
	span, ctx := opentracing.StartSpanFromContext(ctx, "SendPaidNotificationCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := SendPaidNotificationOutput{}

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

	err = c.invoiceService.SendPaidInvoiceNotification(ctx, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionCompleted, result, err // failed send invoice email is not a blocker, since a new attempt will be made automatically by cron
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
