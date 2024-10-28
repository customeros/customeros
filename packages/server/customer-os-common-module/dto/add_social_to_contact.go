package dto

type AddSocialToContact struct {
	SocialId string `json:"socialId"`
	Social   string `json:"social"`
}

func NewAddSocialToContact(socialId, social string) AddSocialToContact {
	output := AddSocialToContact{
		SocialId: socialId,
		Social:   social,
	}
	return output
}
