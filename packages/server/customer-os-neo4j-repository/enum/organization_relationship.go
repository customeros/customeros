package enum

type OrganizationRelationship string

const (
	OrganizationRelationshipProspect       OrganizationRelationship = "PROSPECT"
	OrganizationRelationshipCustomer       OrganizationRelationship = "CUSTOMER"
	OrganizationRelationshipFormerCustomer OrganizationRelationship = "FORMER_CUSTOMER"
	OrganizationRelationshipNotAFit        OrganizationRelationship = "NOT_A_FIT"
)

func (e OrganizationRelationship) String() string {
	return string(e)
}

func DecodeOrganizationRelationship(str string) OrganizationRelationship {
	switch str {
	case OrganizationRelationshipProspect.String():
		return OrganizationRelationshipProspect
	case OrganizationRelationshipCustomer.String():
		return OrganizationRelationshipCustomer
	case OrganizationRelationshipNotAFit.String():
		return OrganizationRelationshipNotAFit
	case OrganizationRelationshipFormerCustomer.String():
		return OrganizationRelationshipFormerCustomer
	default:
		return ""
	}
}

func (e OrganizationRelationship) IsValid() bool {
	switch e {
	case OrganizationRelationshipProspect, OrganizationRelationshipCustomer, OrganizationRelationshipNotAFit, OrganizationRelationshipFormerCustomer:
		return true
	}
	return false
}

func (e OrganizationRelationship) DefaultStage() OrganizationStage {
	switch e {
	case OrganizationRelationshipProspect:
		return Lead
	case OrganizationRelationshipCustomer:
		return Onboarding
	case OrganizationRelationshipNotAFit:
		return Unqualified
	case OrganizationRelationshipFormerCustomer:
		return Target
	default:
		return ""
	}
}
