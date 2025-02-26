package agent_capability

import (
	"context"
	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
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

type GenerateInvoiceConfig struct {
	FromEmail                  ConfigSingleValue     `json:"fromEmail"`
	CcEmails                   ConfigMultipleValues  `json:"ccEmails"`
	BccEmails                  ConfigMultipleValues  `json:"bccEmails"`
	LogoRepositoryFileId       ConfigSingleValue     `json:"logoRepositoryFileId"`
	Country                    ConfigSingleValue     `json:"country"`
	LegalName                  ConfigSingleValue     `json:"legalName"`
	AddressLine1               ConfigSingleValue     `json:"addressLine1"`
	AddressLine2               ConfigSingleValue     `json:"addressLine2"`
	ZIP                        ConfigSingleValue     `json:"zip"`
	Locality                   ConfigSingleValue     `json:"locality"`
	Region                     ConfigSingleValue     `json:"region"`
	IncludeBankTransferDetails ConfigSingleBoolValue `json:"includeBankTransferDetails"`
	BankInfoTemplate           ConfigSingleValue     `json:"bankInfoTemplate"`
	BankName                   ConfigSingleValue     `json:"bankName"`
	AccountNumber              ConfigSingleValue     `json:"accountNumber"`
	IBAN                       ConfigSingleValue     `json:"iban"`
	BIC                        ConfigSingleValue     `json:"bic"`
	SortCode                   ConfigSingleValue     `json:"sortCode"`
	RoutingNumber              ConfigSingleValue     `json:"routingNumber"`
	OtherDetails               ConfigSingleValue     `json:"otherDetails"`
}

func (c *GenerateInvoiceConfig) Validate() bool {
	if c.LegalName.Value == "" {
		c.LegalName.Error = "Please provide your company legal name"
		return false
	} else {
		c.LegalName.Error = ""
	}

	if c.FromEmail.Value == "" {
		c.FromEmail.Error = "Please provide your company email"
		return false
	} else {
		c.FromEmail.Error = ""
	}

	syntaxValidation := mailsherpa.ValidateEmailSyntax(c.FromEmail.Value)
	if !syntaxValidation.IsValid {
		c.FromEmail.Error = "Please provide valid email address"
		return false
	} else {
		c.FromEmail.Error = ""
	}

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
	return "Generate invoices"
}

func (c *GenerateInvoiceCapability) NewInput() GenerateInvoiceInput {
	return GenerateInvoiceInput{}
}

func (c *GenerateInvoiceCapability) NewConfig() GenerateInvoiceConfig {
	return GenerateInvoiceConfig{}
}

func (c *GenerateInvoiceCapability) DefaultConfig() any {
	config := c.NewConfig()
	config.CcEmails.Value = []string{}
	config.BccEmails.Value = []string{}
	return &config
}

func (c *GenerateInvoiceCapability) DefaultActive() bool {
	return true
}

func (c *GenerateInvoiceCapability) ValidateConfig(GenerateInvoiceConfig) error {
	return nil
}

func (c *GenerateInvoiceCapability) ValidateInput(input GenerateInvoiceInput) error {
	if input.ContractId == "" {
		return errors.New("missing required input: ContractId")
	}
	return nil
}

func (c *GenerateInvoiceCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[GenerateInvoiceInput, GenerateInvoiceConfig]) (enum.CapabilityExecutionStatus, GenerateInvoiceOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GenerateInvoiceCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "executionContainer", executionContainer)

	result := GenerateInvoiceOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	dataFields := data_fields.InvoiceFields{
		DryRun:    executionContainer.InputData.DryRun,
		Preview:   executionContainer.InputData.Preview,
		FromEmail: executionContainer.ConfigData.FromEmail.Value,
		CcEmails:  executionContainer.ConfigData.CcEmails.Value,
		BccEmails: executionContainer.ConfigData.BccEmails.Value,
		TenantBillingProfile: &data_fields.TenantBillingProfile{
			LogoRepositoryFileId:       executionContainer.ConfigData.LogoRepositoryFileId.Value,
			LegalName:                  executionContainer.ConfigData.LegalName.Value,
			Country:                    executionContainer.ConfigData.Country.Value,
			Region:                     executionContainer.ConfigData.Region.Value,
			AddressLine1:               executionContainer.ConfigData.AddressLine1.Value,
			AddressLine2:               executionContainer.ConfigData.AddressLine2.Value,
			Zip:                        executionContainer.ConfigData.ZIP.Value,
			Locality:                   executionContainer.ConfigData.Locality.Value,
			IncludeBankTransferDetails: executionContainer.ConfigData.IncludeBankTransferDetails.Value,
			BankName:                   executionContainer.ConfigData.BankName.Value,
			AccountNumber:              executionContainer.ConfigData.AccountNumber.Value,
			IBAN:                       executionContainer.ConfigData.IBAN.Value,
			BIC:                        executionContainer.ConfigData.BIC.Value,
			SortCode:                   executionContainer.ConfigData.SortCode.Value,
			RoutingNumber:              executionContainer.ConfigData.RoutingNumber.Value,
			OtherDetails:               executionContainer.ConfigData.OtherDetails.Value,
		},
	}
	invoiceId, err := c.invoiceService.InvoiceContract(ctx, nil, executionContainer.InputData.ContractId, dataFields)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, result, err
	}
	// load invoice after generation
	invoiceEntity, err := c.invoiceService.GetById(ctx, nil, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get invoice"))
	}
	result.InvoiceId = invoiceId
	if invoiceEntity != nil {
		result.InvoiceNumber = invoiceEntity.Number
	}

	tracing.LogObjectAsJson(span, "result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}
