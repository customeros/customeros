package model

import "strings"

type ReferencedContact struct {
	ExternalId string `json:"externalId,omitempty"`
	Id         string `json:"id,omitempty"`
}

func (r *ReferencedContact) GetReferencedEntityType() ReferencedEntityType {
	return ReferencedEntityTypeContact
}

func (r *ReferencedContact) ReferencedById() bool {
	return strings.TrimSpace(r.Id) != ""
}

func (r *ReferencedContact) ReferencedByExternalId() bool {
	return strings.TrimSpace(r.ExternalId) != "" && strings.TrimSpace(r.Id) == ""
}

func (r *ReferencedContact) Available() bool {
	return r.ReferencedById() || r.ReferencedByExternalId()
}
