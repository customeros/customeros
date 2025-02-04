package agent_capability

import (
	"context"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type GenerateInvoiceCapability struct {
	postgresRepositories *postgres_repository.Repositories
	invoiceService       interfaces.InvoiceService
}

type GenerateInvoiceInput struct {
	ContractId string `json:"contractId"`
	DryRun     bool   `json:"dryRun"`
	Preview    bool   `json:"preview"`
}

type GenerateInvoiceOutput struct {
	InvoiceId     string `json:"invoiceId"`
	InvoiceNumber string `json:"invoiceNumber"`
}

type GenerateInvoiceConfig struct{}

func (c *GenerateInvoiceConfig) Validate() bool {
	return true
}

func NewGenerateInvoiceCapability(postgres *postgres_repository.Repositories, invoiceService interfaces.InvoiceService) *GenerateInvoiceCapability {
	return &GenerateInvoiceCapability{
		postgresRepositories: postgres,
		invoiceService:       invoiceService,
	}
}

// Compile-time interface check
var (
	_ interfaces.AgentCapability[GenerateInvoiceInput, GenerateInvoiceOutput, GenerateInvoiceConfig] = (*GenerateInvoiceCapability)(nil)
)

func (c *GenerateInvoiceCapability) Type() enum.AgentCapability {
	return enum.CapabilityGenerateInvoice
}

func (c *GenerateInvoiceCapability) Name() string {
	return "Generate an invoice"
}

func (c *GenerateInvoiceCapability) NewInput() GenerateInvoiceInput {
	return GenerateInvoiceInput{}
}

func (c *GenerateInvoiceCapability) NewConfig() GenerateInvoiceConfig {
	return GenerateInvoiceConfig{}
}

func (c *GenerateInvoiceCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *GenerateInvoiceCapability) ValidateConfig(config GenerateInvoiceConfig) error {
	return nil
}

func (c *GenerateInvoiceCapability) ValidateInput(input GenerateInvoiceInput) error {
	if input.ContractId == "" {
		return errors.New("missing required input: ContractId")
	}
	return nil
}

func (c *GenerateInvoiceCapability) Execute(ctx context.Context, data GenerateInvoiceInput, config GenerateInvoiceConfig) (GenerateInvoiceOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GenerateInvoiceCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", data)
	tracing.LogObjectAsJson(span, "config", config)

	result := GenerateInvoiceOutput{}

	if err := c.ValidateInput(data); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return result, err
	}
	if err := c.ValidateConfig(config); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return result, err
	}

	tracing.LogObjectAsJson(span, "result", result)
	return result, nil
}
