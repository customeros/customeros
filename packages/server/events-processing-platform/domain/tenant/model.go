package invoice

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"time"
)

type Tenant struct {
	Id              string                 `json:"id"`
	Name            string                 `json:"name"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	SourceFields    events.Source          `json:"source"`
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
	Id                     string        `json:"id"`
	CreatedAt              time.Time     `json:"createdAt"`
	UpdatedAt              time.Time     `json:"updatedAt"`
	SourceFields           events.Source `json:"source"`
	Phone                  string        `json:"phone"`
	AddressLine1           string        `json:"addressLine1"`
	AddressLine2           string        `json:"addressLine2"`
	AddressLine3           string        `json:"addressLine3"`
	Locality               string        `json:"locality"`
	Country                string        `json:"country"`
	Region                 string        `json:"region"`
	Zip                    string        `json:"zip"`
	LegalName              string        `json:"legalName"`
	VatNumber              string        `json:"vatNumber"`
	SendInvoicesFrom       string        `json:"sendInvoicesFrom"`
	SendInvoicesBcc        string        `json:"sendInvoicesBcc"`
	CanPayWithPigeon       bool          `json:"canPayWithPigeon"`
	CanPayWithBankTransfer bool          `json:"canPayWithBankTransfer"`
	Check                  bool          `json:"check"`
}

type BankAccount struct {
	Id                  string        `json:"id"`
	CreatedAt           time.Time     `json:"createdAt"`
	UpdatedAt           time.Time     `json:"updatedAt"`
	SourceFields        events.Source `json:"source"`
	BankName            string        `json:"bankName"`
	BankTransferEnabled bool          `json:"bankTransferEnabled"`
	AllowInternational  bool          `json:"allowInternational"`
	Currency            string        `json:"currency"`
	Iban                string        `json:"iban"`
	Bic                 string        `json:"bic"`
	SortCode            string        `json:"sortCode"`
	AccountNumber       string        `json:"accountNumber"`
	RoutingNumber       string        `json:"routingNumber"`
	OtherDetails        string        `json:"otherDetails"`
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
