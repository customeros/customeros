package entity

import "time"

const (
	Partner  string = "Partner"
	Customer string = "Customer"
	Reseller string = "Reseller"
	Vendor   string = "Vendor"
)

const (
	Prospect string = "Prospect"
	Live     string = "Live"
)

type OrganizationRelation string

const (
	Subsidiary OrganizationRelation = "subsidiary"
)

type OrganizationNote struct {
	FieldSource string
	Note        string
}

type ParentOrganization struct {
	ExternalId           string
	OrganizationRelation OrganizationRelation
	Type                 string
}

type OrganizationData struct {
	Id                  string
	Name                string
	Description         string
	Domains             []string
	Notes               []OrganizationNote
	Website             string
	Industry            string
	IsPublic            bool
	Employees           int64
	PhoneNumber         string
	Email               string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	ExternalId          string
	ExternalSystem      string
	ExternalUrl         string
	ExternalSyncId      string
	UserExternalOwnerId string

	LocationName string
	Country      string
	Region       string
	Locality     string
	Address      string
	Address2     string
	Zip          string

	RelationshipName  string
	RelationshipStage string

	ParentOrganization *ParentOrganization
}

func (o OrganizationData) HasDomains() bool {
	return len(o.Domains) > 0
}

func (o OrganizationData) HasLocation() bool {
	return len(o.LocationName) > 0 || len(o.Country) > 0 || len(o.Region) > 0 || len(o.Locality) > 0 || len(o.Address) > 0 || len(o.Address2) > 0 || len(o.Zip) > 0
}

func (o OrganizationData) HasNotes() bool {
	return len(o.Notes) > 0
}

func (o OrganizationData) HasRelationship() bool {
	return o.RelationshipName != ""
}

func (o OrganizationData) HasPhoneNumber() bool {
	return len(o.PhoneNumber) > 0
}

func (o OrganizationData) HasEmail() bool {
	return len(o.Email) > 0
}

func (o OrganizationData) IsSubsidiary() bool {
	return o.ParentOrganization != nil && o.ParentOrganization.OrganizationRelation == Subsidiary
}

func (o OrganizationData) HasOwner() bool {
	return o.UserExternalOwnerId != ""
}
