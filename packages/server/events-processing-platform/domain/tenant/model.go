package invoice

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type Tenant struct {
	Id              string                 `json:"id"`
	Name            string                 `json:"name"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	SourceFields    commonmodel.Source     `json:"source"`
	BillingProfiles []TenantBillingProfile `json:"billingProfiles"`
	BankAccounts    []BankAccount          `json:"bankAccounts"`
	TenantSettings  TenantSettings         `json:"tenantSettings"`
}

type TenantSettings struct {
	InvoicingEnabled     bool   `json:"invoicingEnabled"`
	InvoicingPostpaid    bool   `json:"invoicingPostpaid"`
	BaseCurrency         string `json:"baseCurrency"`
	LogoUrl              string `json:"logoUrl"`
	LogoRepositoryFileId string `json:"logoRepositoryFileId"`
}

type TenantBillingProfile struct {
	Id                                string             `json:"id"`
	CreatedAt                         time.Time          `json:"createdAt"`
	UpdatedAt                         time.Time          `json:"updatedAt"`
	SourceFields                      commonmodel.Source `json:"source"`
	Phone                             string             `json:"phone"`
	AddressLine1                      string             `json:"addressLine1"`
	AddressLine2                      string             `json:"addressLine2"`
	AddressLine3                      string             `json:"addressLine3"`
	Locality                          string             `json:"locality"`
	Country                           string             `json:"country"`
	Zip                               string             `json:"zip"`
	LegalName                         string             `json:"legalName"`
	DomesticPaymentsBankInfo          string             `json:"domesticPaymentsBankInfo"`
	DomesticPaymentsBankName          string             `json:"domesticPaymentsBankName"`
	DomesticPaymentsAccountNumber     string             `json:"domesticPaymentsAccountNumber"`
	DomesticPaymentsSortCode          string             `json:"domesticPaymentsSortCode"`
	InternationalPaymentsBankInfo     string             `json:"internationalPaymentsBankInfo"`
	InternationalPaymentsSwiftBic     string             `json:"internationalPaymentsSwiftBic"`
	InternationalPaymentsBankName     string             `json:"internationalPaymentsBankName"`
	InternationalPaymentsBankAddress  string             `json:"internationalPaymentsBankAddress"`
	InternationalPaymentsInstructions string             `json:"internationalPaymentsInstructions"`
	VatNumber                         string             `json:"vatNumber"`
	SendInvoicesFrom                  string             `json:"sendInvoicesFrom"`
	SendInvoicesBcc                   string             `json:"sendInvoicesBcc"`
	CanPayWithCard                    bool               `json:"canPayWithCard"`
	CanPayWithDirectDebitSEPA         bool               `json:"canPayWithDirectDebitSEPA"`
	CanPayWithDirectDebitACH          bool               `json:"canPayWithDirectDebitACH"`
	CanPayWithDirectDebitBacs         bool               `json:"canPayWithDirectDebitBacs"`
	CanPayWithPigeon                  bool               `json:"canPayWithPigeon"`
	CanPayWithBankTransfer            bool               `json:"canPayWithBankTransfer"`
}

type BankAccount struct {
	Id                  string             `json:"id"`
	CreatedAt           time.Time          `json:"createdAt"`
	UpdatedAt           time.Time          `json:"updatedAt"`
	SourceFields        commonmodel.Source `json:"source"`
	BankName            string             `json:"bankName"`
	BankTransferEnabled bool               `json:"bankTransferEnabled"`
	AllowInternational  bool               `json:"allowInternational"`
	Currency            string             `json:"currency"`
	Iban                string             `json:"iban"`
	Bic                 string             `json:"bic"`
	SortCode            string             `json:"sortCode"`
	AccountNumber       string             `json:"accountNumber"`
	RoutingNumber       string             `json:"routingNumber"`
	OtherDetails        string             `json:"otherDetails"`
}

func (t Tenant) HasBillingProfile(id string) bool {
	for _, bp := range t.BillingProfiles {
		if bp.Id == id {
			return true
		}
	}
	return false
}

func (t Tenant) GetBillingProfile(id string) *TenantBillingProfile {
	for _, bp := range t.BillingProfiles {
		if bp.Id == id {
			return &bp
		}
	}
	return nil
}

func (t Tenant) HasBankAccount(id string) bool {
	for _, ba := range t.BankAccounts {
		if ba.Id == id {
			return true
		}
	}
	return false
}

func (t Tenant) GetBankAccount(id string) *BankAccount {
	for _, ba := range t.BankAccounts {
		if ba.Id == id {
			return &ba
		}
	}
	return nil
}

func (t Tenant) AddBankAccount(ba BankAccount) {
	t.BankAccounts = append(t.BankAccounts, ba)
}
