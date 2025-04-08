package agent_capability

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type SendInvoiceViaEmailCapability struct {
	postgresRepositories *postgres_repository.Repositories
	invoiceService       interfaces.InvoiceService
}

func NewSendInvoiceViaEmailCapability(postgresRepositories *postgres_repository.Repositories, invoiceService interfaces.InvoiceService) *SendInvoiceViaEmailCapability {
	return &SendInvoiceViaEmailCapability{
		postgresRepositories: postgresRepositories,
		invoiceService:       invoiceService,
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
	return "Send an invoice via email"
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

func (c *SendInvoiceViaEmailCapability) DefaultActive() bool {
	return true
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
	spans, ctx := telemetry.StartServiceSpan(ctx, "SendInvoiceViaEmailCapability.Execute")
	defer spans.Finish()

	spans.LogObjectAsJson("input", executionContainer.InputData)
	spans.LogObjectAsJson("config", executionContainer.ConfigData)

	result := SendInvoiceViaEmailOutput{}

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

	paymentEnabled, err := c.isPaymentCapabilityEnabled(ctx, executionContainer.AgentExecutionID)
	if err != nil {
		spans.TraceError(err)
	}

	if paymentEnabled {
		err = c.invoiceService.GenerateNewPaymentLink(ctx, executionContainer.InputData.InvoiceID)
		if err != nil {
			spans.TraceError(err)
			// Do not stop capability execution if payment link generation fails
		}
	}

	err = c.invoiceService.SendInvoiceNotification(ctx, executionContainer.InputData.InvoiceID, paymentEnabled)
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionCompleted, result, err // failed send invoice email is not a blocker, since a new attempt will be made automatically by cron
	}

	spans.LogObjectAsJson("result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}

func (c *SendInvoiceViaEmailCapability) isPaymentCapabilityEnabled(ctx context.Context, agentExecutionId string) (bool, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "SendInvoiceViaEmailCapability.isPaymentCapabilityEnabled")
	defer spans.Finish()
	spans.LogKV("agentExecutionId", agentExecutionId)

	// get agent id from execution
	agentExecution, err := c.postgresRepositories.AgentExecutionRepository.GetById(ctx, agentExecutionId)
	if err != nil {
		spans.TraceError(err)
		return false, err
	}
	if agentExecution == nil || utils.IfNotNilString(agentExecution.AgentID) == "" {
		err := errors.New("agent execution not found")
		spans.TraceError(err)
		return false, err
	}
	agentId := utils.IfNotNilString(agentExecution.AgentID)

	// get agent
	agent, err := c.postgresRepositories.AgentRepository.GetById(ctx, agentId)
	if err != nil {
		spans.TraceError(err)
		return false, err
	}
	if agent == nil {
		err := errors.New("agent not found")
		spans.TraceError(err)
		return false, err
	}

	// check capability
	for _, capability := range agent.Capabilities {
		if capability.Type == enum.CapabilityProcessAutopayment {
			return capability.Active, nil
		}
	}

	return false, nil
}
