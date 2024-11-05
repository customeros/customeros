package dto

type SocialFields struct {
	Url           string `json:"url"`
	Alias         string `json:"alias"`
	ExternalId    string `json:"externalId"`
	FollowerCount int    `json:"followerCount"`
}

type CreateSocial struct {
	SocialFields
}
