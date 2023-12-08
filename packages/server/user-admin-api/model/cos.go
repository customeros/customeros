package model

type GetPlayerResponse struct {
	PlayerByAuthIdProvider struct {
		Id    string `json:"id"`
		Users *[]struct {
			Tenant string `json:"tenant"`
		} `json:"users"`
	} `json:"player_ByAuthIdProvider"`
}

type Result struct {
	Result bool `json:"result"`
}

type GetTenantByWorkspaceResponse struct {
	Name *string `json:"tenant_ByWorkspace"`
}

type MergeTenantToWorkspaceResponse struct {
	Workspace struct {
		Result bool `json:"result"`
	} `json:"workspace_MergeToTenant"`
}

type CreateUserResponse struct {
	User struct {
		ID    string  `json:"id"`
		Roles *[]Role `json:"roles"`
	} `json:"user_Create"`
}

type UserAddRoleResponse struct {
	UserAddRole struct {
		ID    string  `json:"id"`
		Roles *[]Role `json:"roles"`
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

type CreateContractResponse struct {
	ContractCreate struct {
		Id string `json:"id"`
	} `json:"contract_Create"`
}

type CreateServiceLineItemResponse struct {
	ServiceLineItemCreate struct {
		Id string `json:"id"`
	} `json:"serviceLineItemCreate"`
}

type CreateMeetingResponse struct {
	MeetingCreate struct {
		Id string `json:"id"`
	} `json:"meeting_Create"`
}

type Role string

const (
	RoleAdmin                   Role = "ADMIN"
	RoleCustomerOsPlatformOwner Role = "CUSTOMER_OS_PLATFORM_OWNER"
	RoleOwner                   Role = "OWNER"
	RoleUser                    Role = "USER"
)

type GetUserByEmailResponse struct {
	UserByEmail struct {
		ID    string  `json:"id"`
		Roles *[]Role `json:"roles"`
	} `json:"user_ByEmail"`
}
type GetUserByIdResponse struct {
	User struct {
		ID    string  `json:"id"`
		Roles *[]Role `json:"roles"`
	} `json:"user"`
}

type UserResponse struct {
	ID    string  `json:"id"`
	Roles *[]Role `json:"roles"`
}
