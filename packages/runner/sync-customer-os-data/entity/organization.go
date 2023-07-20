package entity

import (
	utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type OrganizationRelation string

const (
	Subsidiary OrganizationRelation = "subsidiary"
)

type OrganizationNote struct {
	FieldSource string `json:"fieldSource"`
	Note        string `json:"note"`
}

type ParentOrganization struct {
	ExternalId           string               `json:"externalId,omitempty"`
	OrganizationRelation OrganizationRelation `json:"organizationRelation,omitempty"`
	Type                 string               `json:"type,omitempty"`
}

type OrganizationData struct {
	BaseData
	Name                string              `json:"name,omitempty"`
	Description         string              `json:"description,omitempty"`
	Domains             []string            `json:"domains,omitempty"`
	Notes               []OrganizationNote  `json:"notes,omitempty"`
	Website             string              `json:"website,omitempty"`
	Industry            string              `json:"industry,omitempty"`
	IsPublic            bool                `json:"isPublic,omitempty"`
	Employees           int64               `json:"employees,omitempty"`
	PhoneNumber         string              `json:"phoneNumber,omitempty"`
	Email               string              `json:"email,omitempty"`
	ExternalUrl         string              `json:"externalUrl,omitempty"`
	ExternalSourceTable *string             `json:"externalSourceTable,omitempty"`
	UserExternalOwnerId string              `json:"externalOwnerId,omitempty"`
	UserExternalId      string              `json:"externalUserId,omitempty"`
	LocationName        string              `json:"locationName,omitempty"`
	Country             string              `json:"country,omitempty"`
	Region              string              `json:"region,omitempty"`
	Locality            string              `json:"locality,omitempty"`
	Address             string              `json:"address,omitempty"`
	Address2            string              `json:"address2,omitempty"`
	Zip                 string              `json:"zip,omitempty"`
	RelationshipName    string              `json:"relationshipName,omitempty"`
	RelationshipStage   string              `json:"relationshipStage,omitempty"`
	ParentOrganization  *ParentOrganization `json:"parentOrganization,omitempty"`
	SubIndustry         string              `json:"subIndustry,omitempty"`
	IndustryGroup       string              `json:"industryGroup,omitempty"`
	TargetAudience      string              `json:"targetAudience,omitempty"`
	ValueProposition    string              `json:"valueProposition,omitempty"`
	Market              string              `json:"market,omitempty"`
	LastFundingRound    string              `json:"lastFundingRound,omitempty"`
	LastFundingAmount   string              `json:"lastFundingAmount,omitempty"`
}

func (o *OrganizationData) HasDomains() bool {
	return len(o.Domains) > 0
}

func (o *OrganizationData) HasLocation() bool {
	return len(o.LocationName) > 0 || len(o.Country) > 0 || len(o.Region) > 0 || len(o.Locality) > 0 || len(o.Address) > 0 || len(o.Address2) > 0 || len(o.Zip) > 0
}

func (o *OrganizationData) HasNotes() bool {
	return len(o.Notes) > 0
}

func (o *OrganizationData) HasRelationship() bool {
	return o.RelationshipName != ""
}

func (o *OrganizationData) HasPhoneNumber() bool {
	return len(o.PhoneNumber) > 0
}

func (o *OrganizationData) HasEmail() bool {
	return len(o.Email) > 0
}

func (o *OrganizationData) IsSubsidiary() bool {
	return o.ParentOrganization != nil && o.ParentOrganization.OrganizationRelation == Subsidiary
}

func (o *OrganizationData) HasOwnerByOwnerId() bool {
	return o.UserExternalOwnerId != ""
}

func (o *OrganizationData) HasOwnerByUserId() bool {
	return o.UserExternalId != ""
}

func (o *OrganizationData) Normalize() {
	o.SetTimes()

	utils.FilterEmpty(o.Domains)
	utils.LowercaseStrings(o.Domains)
	o.Domains = utils.RemoveDuplicates(o.Domains)
}
