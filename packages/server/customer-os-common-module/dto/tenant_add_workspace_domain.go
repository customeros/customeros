package dto

type AddWorkspaceDomainToTenant struct {
	Domain   string `json:"domain"`
	Provider string `json:"provider"`
}
