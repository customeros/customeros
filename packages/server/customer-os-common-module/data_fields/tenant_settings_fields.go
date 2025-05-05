package data_fields

// Nil fields wil be skipped from update
type TenantSettingsFields struct {
	InvoicingPostpaid *bool   `json:"invoicingPostpaid,omitempty"`
	BaseCurrency      *string `json:"baseCurrency,omitempty"`
	WorkspaceLogoKey  *string `json:"workspaceLogoKey,omitempty"`
	WorkspaceName     *string `json:"workspaceName,omitempty"`
	WorkspaceLogo     *string `json:"workspaceLogo,omitempty"` // Deprecated
}

// IsEmpty returns true if there are no fields to update
func (fields TenantSettingsFields) IsEmpty() bool {
	return fields.InvoicingPostpaid == nil &&
		fields.BaseCurrency == nil &&
		fields.WorkspaceLogo == nil &&
		fields.WorkspaceLogoKey == nil &&
		fields.WorkspaceName == nil
}
