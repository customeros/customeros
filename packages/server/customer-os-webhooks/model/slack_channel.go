package model

type SlackChannelData struct {
	BaseData
	ChannelId string `json:"channelId,omitempty"`
}

func (u *SlackChannelData) Normalize() {
	u.SetTimes()
	u.BaseData.Normalize()
}
