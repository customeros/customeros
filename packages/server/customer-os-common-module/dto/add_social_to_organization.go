package dto

type AddSocialToOrganization struct {
	SocialId string `json:"socialId"`
	Social   string `json:"social"`
}

func NewAddSocialToOrganization(socialId, social string) AddSocialToOrganization {
	output := AddSocialToOrganization{
		SocialId: socialId,
		Social:   social,
	}
	return output
}
