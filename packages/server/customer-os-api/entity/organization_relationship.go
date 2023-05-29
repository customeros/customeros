package entity

type OrganizationRelationship string

const (
	Investor    OrganizationRelationship = "INVESTOR"
	Supplier    OrganizationRelationship = "SUPPLIER"
	Partner     OrganizationRelationship = "PARTNER"
	Customer    OrganizationRelationship = "CUSTOMER"
	Distributor OrganizationRelationship = "DISTRIBUTOR"
)

func (o OrganizationRelationship) FromString(input string) OrganizationRelationship {
	switch input {
	case "INVESTOR":
		return Investor
	case "SUPPLIER":
		return Supplier
	case "PARTNER":
		return Partner
	case "CUSTOMER":
		return Customer
	case "DISTRIBUTOR":
		return Distributor
	default:
		return ""
	}
}

func (o OrganizationRelationship) String() string {
	return string(o)
}
