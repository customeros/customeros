package model

type GetTenantByWorkspaceResponse struct {
	Workspace *string `json:"tenant_ByWorkspace"`
}

type MergeTenantToWorkspaceResponse struct {
	Workspace struct {
		Result bool `json:"result"`
	} `json:"workspace_MergeToTenant"`
}

type CreateUserResponse struct {
	User struct {
		ID        string `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	} `json:"user_Create"`
}

type CreateTenantResponse struct {
	Tenant struct {
		Result bool `json:"result"`
	} `json:"tenant_Merge"`
}
