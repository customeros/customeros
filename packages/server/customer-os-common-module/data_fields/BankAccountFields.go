package data_fields

// Nil fields wil be skipped from update
type BankAccountFields struct {
	ID                  string  `json:"id,omitempty"`
	AppSource           *string `json:"appSource,omitempty"`
	Source              *string `json:"source,omitempty"`
	BankName            *string `json:"bankName,omitempty"`
	BankTransferEnabled *bool   `json:"bankTransferEnabled,omitempty"`
	AllowInternational  *bool   `json:"allowInternational,omitempty"`
	Currency            *string `json:"currency,omitempty"`
	IBAN                *string `json:"iban,omitempty"`
	BIC                 *string `json:"bic,omitempty"`
	SortCode            *string `json:"sortCode,omitempty"`
	AccountNumber       *string `json:"accountNumber,omitempty"`
	RoutingNumber       *string `json:"routingNumber,omitempty"`
	OtherDetails        *string `json:"otherDetails,omitempty"`
}

// IsEmpty returns true if there are no fields to update
func (fields BankAccountFields) IsEmpty() bool {
	return fields.AppSource == nil &&
		fields.Source == nil &&
		fields.BankName == nil &&
		fields.BankTransferEnabled == nil &&
		fields.AllowInternational == nil &&
		fields.Currency == nil &&
		fields.IBAN == nil &&
		fields.BIC == nil &&
		fields.SortCode == nil &&
		fields.AccountNumber == nil &&
		fields.RoutingNumber == nil &&
		fields.OtherDetails == nil
}
