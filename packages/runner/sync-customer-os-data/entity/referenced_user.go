package entity

type ReferencedUser struct {
	ExternalId      string `json:"externalId,omitempty"`
	ExternalOwnerId string `json:"externalOwnerId,omitempty"`
	Id              string `json:"id,omitempty"`
}

func (r *ReferencedUser) ReferencedById() bool {
	return r.Id != ""
}

func (r *ReferencedUser) ReferencedByExternalId() bool {
	return r.ExternalId != "" && r.Id == ""
}

func (r *ReferencedUser) ReferencedByExternalOwnerId() bool {
	return r.ExternalOwnerId != "" && r.Id == "" && r.ExternalId == ""
}

func (r *ReferencedUser) Available() bool {
	return r.ReferencedById() || r.ReferencedByExternalId() || r.ReferencedByExternalOwnerId()
}
