package dto

type UpdateDomain struct {
	Primary       bool   `json:"primary"`
	PrimaryDomain string `json:"primary_domain"`
}
