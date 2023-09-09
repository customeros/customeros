package entity

type ReferencedUser struct {
	ExternalId      string `json:"externalId,omitempty"`
	ExternalOwnerId string `json:"externalOwnerId,omitempty"`
	Id              string `json:"id,omitempty"`
}
