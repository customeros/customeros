package data_fields

// Nil fields wil be skipped from update
type TenantSettingsFields struct {
	InvoicingEnabled     *bool   `json:"invoicingEnabled,omitempty"`
	InvoicingPostpaid    *bool   `json:"invoicingPostpaid,omitempty"`
	LogoRepositoryFileId *string `json:"logoRepositoryFileId,omitempty"`
	BaseCurrency         *string `json:"baseCurrency,omitempty"`
	WorkspaceLogo        *string `json:"workspaceLogo,omitempty"`
	WorkspaceName        *string `json:"workspaceName,omitempty"`
}

// IsEmpty returns true if there are no fields to update
func (fields TenantSettingsFields) IsEmpty() bool {
	return fields.InvoicingEnabled == nil &&
		fields.InvoicingPostpaid == nil &&
		fields.LogoRepositoryFileId == nil &&
		fields.BaseCurrency == nil &&
		fields.WorkspaceLogo == nil &&
		fields.WorkspaceName == nil
}
