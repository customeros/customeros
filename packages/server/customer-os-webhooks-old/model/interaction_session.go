package model

type InteractionSessionData struct {
	BaseData
	Name        string `json:"name,omitempty"`
	Channel     string `json:"channel,omitempty"`
	ChannelData string `json:"channelData,omitempty"`
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	Identifier  string `json:"identifier,omitempty"`
}

func (i *InteractionSessionData) Normalize() {
	i.SetTimes()
	i.BaseData.Normalize()
}
