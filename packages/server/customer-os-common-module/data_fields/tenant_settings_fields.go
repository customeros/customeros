package data_fields

// Nil fields wil be skipped from update
type TenantSettingsFields struct {
	InvoicingPostpaid *bool   `json:"invoicingPostpaid,omitempty"`
	BaseCurrency      *string `json:"baseCurrency,omitempty"`
	WorkspaceLogoUrl  *string `json:"workspaceLogoUrl,omitempty"`
	WorkspaceName     *string `json:"workspaceName,omitempty"`
	WorkspaceLogo     *string `json:"workspaceLogo,omitempty"` // Deprecated
}

// IsEmpty returns true if there are no fields to update
func (fields TenantSettingsFields) IsEmpty() bool {
	return fields.InvoicingPostpaid == nil &&
		fields.BaseCurrency == nil &&
		fields.WorkspaceLogo == nil &&
		fields.WorkspaceLogoUrl == nil &&
		fields.WorkspaceName == nil
}
