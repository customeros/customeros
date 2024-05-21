package model

type GetPlayerResponse struct {
	PlayerByAuthIdProvider struct {
		Id    string `json:"id"`
		Users *[]struct {
			Tenant string `json:"tenant"`
		} `json:"users"`
	} `json:"player_ByAuthIdProvider"`
}

type TenantHardDeleteResponse struct {
	Data []struct {
		TenantHardDelete bool `json:"tenant_hardDelete"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
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

type TenantAddBillingProfileResponse struct {
	TenantBillingProfileAdd struct {
		Id string `json:"id"`
	} `json:"tenant_AddBillingProfile"`
}

type NextInvoiceDryRunForContractResponse struct {
	Invoice struct {
		Id string `json:"id"`
	} `json:"invoice_NextDryRunForContract"`
}

type ArchiveOrganizationResponse struct {
	Result bool `json:"result"`
}

type CreateOrganizationResponse struct {
	OrganizationCreate struct {
		Id string `json:"id"`
	} `json:"organization_Create"`
}

type UpdateOrganizationResponse struct {
	OrganizationUpdate struct {
		Id string `json:"id"`
	} `json:"organization_Update"`
}

type GetOrganizationsResponse struct {
	Organizations struct {
		Content       []Organization `json:"content"`
		TotalElements int64          `json:"totalElements"`
	} `json:"organizations"`
}

type Organization struct {
	ID string `json:"id"`
}

type CreateContractResponse struct {
	ContractCreate struct {
		Id string `json:"id"`
	} `json:"contract_Create"`
}

type UpdateContractResponse struct {
	ContractUpdate struct {
		Id string `json:"id"`
	} `json:"contract_Update"`
}

type CreateServiceLineItemResponse struct {
	ContractLineItemCreate struct {
		Metadata struct {
			Id string `json:"id"`
		} `json:"metadata"`
	} `json:"contractLineItem_Create"`
}

type CreateMasterPlanResponse struct {
	MasterPlanCreate struct {
		Id string `json:"id"`
	} `json:"masterPlan_Create"`
}
type CreateMasterPlanMilestoneResponse struct {
	MasterPlanMilestoneCreate struct {
		Id string `json:"id"`
	} `json:"masterPlanMilestone_Create"`
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
	RolePlatformOwner           Role = "PLATFORM_OWNER"
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

type CreateExternalSystemResponse struct {
	Id string `json:"externalSystem_Create"`
}
