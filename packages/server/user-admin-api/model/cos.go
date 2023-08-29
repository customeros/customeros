package model

type GetPlayerResponse struct {
	Id string `json:"id"`
}

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

type AddRoleResponse struct {
	Role struct {
		ID string `json:"id"`
	} `json:"user_AddRole"`
}

type CreateTenantResponse struct {
	Tenant string `json:"tenant_Merge"`
}

type CreateOrganizationResponse struct {
	OrganizationCreate struct {
		Id string `json:"id"`
	} `json:"organization_Create"`
}

type CreateMeetingResponse struct {
	MeetingCreate struct {
		Id string `json:"id"`
	} `json:"meeting_Create"`
}
