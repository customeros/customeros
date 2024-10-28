package dto

type CreateSocial struct {
	Url           string `json:"url"`
	Alias         string `json:"alias"`
	ExtId         string `json:"externalId"`
	FollowerCount int64  `json:"followerCount"`
}
