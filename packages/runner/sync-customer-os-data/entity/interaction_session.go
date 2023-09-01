package entity

type InteractionSession struct {
	BaseData
	Name       string `json:"name,omitempty"`
	Channel    string `json:"channel,omitempty"`
	Type       string `json:"type,omitempty"`
	Status     string `json:"status,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}
