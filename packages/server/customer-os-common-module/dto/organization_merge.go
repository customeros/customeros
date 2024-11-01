package dto

type MergeOrganizations struct {
	SourceOrgId string `json:"sourceOrganizationId"`
	TargetOrgId string `json:"targetOrganizationId"`
}
