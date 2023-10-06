package model

type ReferencedUser struct {
	ExternalId       string `json:"externalId,omitempty"`
	ExternalIdSecond string `json:"externalIdSecond,omitempty"`
	Id               string `json:"id,omitempty"`
}

func (r *ReferencedUser) GetReferencedEntityType() ReferencedEntityType {
	return ReferencedEntityTypeUser
}

func (r *ReferencedUser) ReferencedById() bool {
	return r.Id != ""
}

func (r *ReferencedUser) ReferencedByExternalId() bool {
	return r.ExternalId != "" && r.Id == ""
}

func (r *ReferencedUser) ReferencedByExternalOwnerId() bool {
	return r.ExternalIdSecond != "" && r.Id == "" && r.ExternalId == ""
}

func (r *ReferencedUser) Available() bool {
	return r.ReferencedById() || r.ReferencedByExternalId() || r.ReferencedByExternalOwnerId()
}
