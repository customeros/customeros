package enum

type OrganizationRelationship string

const (
	Prospect       OrganizationRelationship = "PROSPECT"
	Customer       OrganizationRelationship = "CUSTOMER"
	Stranger       OrganizationRelationship = "STRANGER"
	FormerCustomer OrganizationRelationship = "FORMER_CUSTOMER"
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
	case Stranger.String():
		return Stranger
	case FormerCustomer.String():
		return FormerCustomer
	default:
		return ""
	}
}

func (e OrganizationRelationship) IsValid() bool {
	switch e {
	case Prospect, Customer, Stranger, FormerCustomer:
		return true
	}
	return false
}
