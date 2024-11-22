package model

type ReferencedInteractionSession struct {
	ExternalId string `json:"externalId,omitempty"`
}

func (r *ReferencedInteractionSession) Available() bool {
	return r.ReferencedByExternalId()
}

func (r *ReferencedInteractionSession) GetReferencedEntityType() ReferencedEntityType {
	return ReferencedEntityTypeSession
}

func (r *ReferencedInteractionSession) ReferencedByExternalId() bool {
	return r.ExternalId != ""
}
