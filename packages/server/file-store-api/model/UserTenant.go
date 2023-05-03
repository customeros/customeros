package model

type UserTenant struct {
	UserEmail  string `json:"userEmail,omitempty"`
	TenantName string `json:"tenantName,omitempty"`
	UserId     string `json:"userId,omitempty"`
}
