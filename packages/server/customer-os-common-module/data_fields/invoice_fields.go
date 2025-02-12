package data_fields

import (
	"time"
)

type InvoiceFields struct {
	DryRun               bool                  `json:"dryRun"`
	Preview              bool                  `json:"preview"`
	InvoiceStartDate     *time.Time            `json:"invoiceStartDate"`
	InvoiceEndDate       *time.Time            `json:"invoiceEndDate"`
	InvoiceNumber        string                `json:"invoiceNumber"`
	TenantBillingProfile *TenantBillingProfile `json:"tenantBillingProfile"`
}

type TenantBillingProfile struct {
	Country                    string `json:"country"`
	LegalName                  string `json:"legalName"`
	AddressLine1               string `json:"addressLine1"`
	AddressLine2               string `json:"addressLine2"`
	Zip                        string `json:"zip"`
	Locality                   string `json:"locality"`
	Region                     string `json:"region"`
	IncludeBankTransferDetails bool   `json:"includeBankTransfer"`
	BankName                   string `json:"bankName"`
	AccountNumber              string `json:"accountNumber"`
	IBAN                       string `json:"iban"`
	BIC                        string `json:"bic"`
	SortCode                   string `json:"sortCode"`
	RoutingNumber              string `json:"routingNumber"`
	OtherDetails               string `json:"otherDetails"`
}

// Generate dummy, aka acme data
func (f *InvoiceFields) FillWithDummyTenantBillingProfile() {
	f.TenantBillingProfile = &TenantBillingProfile{}
	f.TenantBillingProfile.Country = "US"
	f.TenantBillingProfile.LegalName = "Acme Inc."
	f.TenantBillingProfile.AddressLine1 = "123 Main St."
	f.TenantBillingProfile.AddressLine2 = "Suite 100"
	f.TenantBillingProfile.Zip = "12345"
	f.TenantBillingProfile.Locality = "Anytown"
	f.TenantBillingProfile.Region = "CA"
	f.TenantBillingProfile.IncludeBankTransferDetails = true
	f.TenantBillingProfile.BankName = "Acme Bank"
	f.TenantBillingProfile.AccountNumber = "1234567890"
	f.TenantBillingProfile.IBAN = "US1234567890"
	f.TenantBillingProfile.BIC = "ACMEUS"
	f.TenantBillingProfile.SortCode = "123456"
	f.TenantBillingProfile.RoutingNumber = "123456789"
	f.TenantBillingProfile.OtherDetails = "Invoice number in payment reference"
}
