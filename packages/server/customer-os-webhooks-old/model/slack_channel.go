package model

type SlackChannelData struct {
	BaseData
	ChannelId   string `json:"channelId,omitempty"`
	ChannelName string `json:"channelName,omitempty"`
}

func (u *SlackChannelData) Normalize() {
	u.SetTimes()
	u.BaseData.Normalize()
}
