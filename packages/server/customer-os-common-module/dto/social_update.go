package dto

type UpdateSocial struct {
	Url           string `json:"url"`
	Alias         string `json:"alias"`
	ExternalId    string `json:"externalId"`
	FollowerCount int    `json:"followerCount"`
}
