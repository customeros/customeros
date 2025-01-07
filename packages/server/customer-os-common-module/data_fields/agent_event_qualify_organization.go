package data_fields

import ()

type OrganizationQualifyEventFields struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
}

func (f OrganizationQualifyEventFields) Type() string {
	return "OrganizationQualifyEventFields"
}
