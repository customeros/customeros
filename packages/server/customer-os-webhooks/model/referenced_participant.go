package model

type ReferencedParticipant struct {
	ExternalId string `json:"externalId,omitempty"`
}

func (r *ReferencedParticipant) Available() bool {
	return r.ExternalId != ""
}

func (r *ReferencedParticipant) GetReferencedEntityType() ReferencedEntityType {
	return ReferencedEntityTypeUnknown
}
