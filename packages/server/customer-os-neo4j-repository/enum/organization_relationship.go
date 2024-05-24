package enum

type OrganizationRelationship string

const (
	Prospect       OrganizationRelationship = "PROSPECT"
	Customer       OrganizationRelationship = "CUSTOMER"
	FormerCustomer OrganizationRelationship = "FORMER_CUSTOMER"
	NotAFit        OrganizationRelationship = "NOT_A_FIT"
)

func (e OrganizationRelationship) String() string {
	return string(e)
}

func DecodeOrganizationRelationship(str string) OrganizationRelationship {
	switch str {
	case Prospect.String():
		return Prospect
	case Customer.String():
		return Customer
	case NotAFit.String():
		return NotAFit
	case FormerCustomer.String():
		return FormerCustomer
	default:
		return ""
	}
}

func (e OrganizationRelationship) IsValid() bool {
	switch e {
	case Prospect, Customer, NotAFit, FormerCustomer:
		return true
	}
	return false
}
