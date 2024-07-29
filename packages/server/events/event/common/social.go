package common

type Social struct {
	Url            string `json:"url"`
	Alias          string `json:"alias"`
	ExternalId     string `json:"externalId"`
	FollowersCount int64  `json:"followersCount"`
}
