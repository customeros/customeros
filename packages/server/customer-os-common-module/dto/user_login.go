package dto

type UserLogin struct {
	LoginEmail string `json:"loginEmail"`
	Provider   string `json:"provider"`
	IdentityId string `json:"identityId"`
}
