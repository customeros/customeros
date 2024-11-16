package data_fields

// Nil fields wil be skipped from update
type TenantBillingProfileFields struct {
	ID                     string  `json:"id,omitempty"`
	AppSource              *string `json:"appSource,omitempty"`
	Source                 *string `json:"source,omitempty"`
	Phone                  *string `json:"phone,omitempty"`
	LegalName              *string `json:"legalName,omitempty"`
	AddressLine1           *string `json:"addressLine1,omitempty"`
	AddressLine2           *string `json:"addressLine2,omitempty"`
	AddressLine3           *string `json:"addressLine3,omitempty"`
	Locality               *string `json:"locality,omitempty"`
	Country                *string `json:"country,omitempty"`
	Region                 *string `json:"region,omitempty"`
	Zip                    *string `json:"zip,omitempty"`
	VatNumber              *string `json:"vatNumber,omitempty"`
	SendInvoicesFrom       *string `json:"sendInvoicesFrom,omitempty"`
	SendInvoicesBcc        *string `json:"sendInvoicesBcc,omitempty"`
	CanPayWithPigeon       *bool   `json:"canPayWithPigeon,omitempty"`
	CanPayWithBankTransfer *bool   `json:"canPayWithBankTransfer,omitempty"`
	Check                  *bool   `json:"check,omitempty"`
}

// IsEmpty returns true if there are no fields to update
func (fields TenantBillingProfileFields) IsEmpty() bool {
	return fields.Phone == nil &&
		fields.LegalName == nil &&
		fields.AddressLine1 == nil &&
		fields.AddressLine2 == nil &&
		fields.AddressLine3 == nil &&
		fields.Locality == nil &&
		fields.Country == nil &&
		fields.Region == nil &&
		fields.Zip == nil &&
		fields.VatNumber == nil &&
		fields.SendInvoicesFrom == nil &&
		fields.SendInvoicesBcc == nil &&
		fields.CanPayWithPigeon == nil &&
		fields.CanPayWithBankTransfer == nil &&
		fields.Check == nil
}
