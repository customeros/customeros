package entity

type ReferencedJobRole struct {
	ReferencedContact      ReferencedContact      `json:"referencedContact,omitempty"`
	ReferencedOrganization ReferencedOrganization `json:"referencedOrganization,omitempty"`
}

func (r *ReferencedJobRole) Available() bool {
	return r.ReferencedContact.Available() && r.ReferencedOrganization.Available()
}
