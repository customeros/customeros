package entity

type ReferencedParticipant struct {
	ExternalId string `json:"externalId,omitempty"`
}

func (r *ReferencedParticipant) Available() bool {
	return r.ExternalId != ""
}
