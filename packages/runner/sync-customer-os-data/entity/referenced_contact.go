package entity

type ReferencedContact struct {
	ExternalId string `json:"externalId,omitempty"`
	Id         string `json:"id,omitempty"`
}

func (r *ReferencedContact) GetReferencedEntityType() ReferencedEntityType {
	return ReferencedEntityTypeContact
}

func (r *ReferencedContact) ReferencedById() bool {
	return r.Id != ""
}

func (r *ReferencedContact) ReferencedByExternalId() bool {
	return r.ExternalId != "" && r.Id == ""
}

func (r *ReferencedContact) Available() bool {
	return r.ReferencedById() || r.ReferencedByExternalId()
}
