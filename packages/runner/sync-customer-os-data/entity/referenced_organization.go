package entity

type ReferencedOrganization struct {
	ExternalId string `json:"externalId,omitempty"`
	Id         string `json:"id,omitempty"`
	Domain     string `json:"domain,omitempty"`
	JobTitle   string `json:"jobTitle,omitempty"`
}

func (r *ReferencedOrganization) GetReferencedEntityType() ReferencedEntityType {
	return ReferencedEntityTypeOrganization
}

func (r *ReferencedOrganization) ReferencedById() bool {
	return r.Id != ""
}

func (r *ReferencedOrganization) ReferencedByExternalId() bool {
	return r.ExternalId != "" && r.Id == ""
}

func (r *ReferencedOrganization) ReferencedByDomain() bool {
	return r.Domain != "" && r.Id == "" && r.ExternalId == ""
}

func (r *ReferencedOrganization) Available() bool {
	return r.ReferencedById() || r.ReferencedByExternalId() || r.ReferencedByDomain()
}
