package data_fields

// Nil fields wil be skipped from update
type TenantSettingsFields struct {
	InvoicingPostpaid       *bool   `json:"invoicingPostpaid,omitempty"`
	BaseCurrency            *string `json:"baseCurrency,omitempty"`
	WorkspaceLogo           *string `json:"workspaceLogo,omitempty"` // Deprecated
	WorkspaceLogoIdentifier *string `json:"workspaceLogoIdentifier,omitempty"`
	WorkspaceName           *string `json:"workspaceName,omitempty"`
}

// IsEmpty returns true if there are no fields to update
func (fields TenantSettingsFields) IsEmpty() bool {
	return fields.InvoicingPostpaid == nil &&
		fields.BaseCurrency == nil &&
		fields.WorkspaceLogo == nil &&
		fields.WorkspaceLogoIdentifier == nil &&
		fields.WorkspaceName == nil
}
