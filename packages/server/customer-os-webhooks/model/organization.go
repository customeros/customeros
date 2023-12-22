package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
)

type OrganizationNote struct {
	FieldSource string `json:"fieldSource"`
	Note        string `json:"note"`
}

type ParentOrganization struct {
	Organization ReferencedOrganization `json:"organization,omitempty"`
	Type         string                 `json:"type,omitempty"`
}

type OrganizationData struct {
	BaseData
	CustomerOsId string `json:"customerOsId"`
	Name         string `json:"name,omitempty"`
	// Fallback name is used when name is empty and domain is missing
	FallbackName string `json:"fallbackName,omitempty"`
	Description  string `json:"description,omitempty"`
	// For sub-orgs this property is ignored
	Domains      []string           `json:"domains,omitempty"`
	Notes        []OrganizationNote `json:"notes,omitempty"` // Deprecated, decide what to do with multi notes
	Note         string             `json:"note,omitempty"`
	ReferenceId  string             `json:"referenceId,omitempty"`
	Website      string             `json:"website,omitempty"`
	Industry     string             `json:"industry,omitempty"`
	IsPublic     bool               `json:"isPublic,omitempty"`
	IsCustomer   bool               `json:"isCustomer,omitempty"`
	Employees    int64              `json:"employees,omitempty"`
	PhoneNumbers []PhoneNumber      `json:"phoneNumbers,omitempty"`
	Email        string             `json:"email,omitempty"`
	// Currently not used. Sync processes will not set automatically owner user
	OwnerUser          *ReferencedUser     `json:"ownerUser,omitempty"`
	LocationName       string              `json:"locationName,omitempty"`
	Country            string              `json:"country,omitempty"`
	Region             string              `json:"region,omitempty"`
	Locality           string              `json:"locality,omitempty"`
	Address            string              `json:"address,omitempty"`
	Address2           string              `json:"address2,omitempty"`
	Zip                string              `json:"zip,omitempty"`
	ParentOrganization *ParentOrganization `json:"parentOrganization,omitempty"`
	SubIndustry        string              `json:"subIndustry,omitempty"`
	IndustryGroup      string              `json:"industryGroup,omitempty"`
	TargetAudience     string              `json:"targetAudience,omitempty"`
	ValueProposition   string              `json:"valueProposition,omitempty"`
	Market             string              `json:"market,omitempty"`
	LastFundingRound   string              `json:"lastFundingRound,omitempty"`
	LastFundingAmount  string              `json:"lastFundingAmount,omitempty"`
	LogoUrl            string              `json:"logoUrl,omitempty"`
	YearFounded        *int64              `json:"yearFounded,omitempty"`
	Headquarters       string              `json:"headquarters,omitempty"`
	EmployeeGrowthRate string              `json:"employeeGrowthRate,omitempty"`
	// If true, the organization will be created by domain,
	// Missing domains, or blacklisted domains will result in no organization being created
	// Note: for sub-orgs this property is ignored
	DomainRequired bool `json:"domainRequired"`
	Whitelisted    bool `json:"whitelisted"`
	UpdateOnly     bool `json:"updateOnly"`
}

func (o *OrganizationData) HasDomains() bool {
	return len(o.Domains) > 0
}

func (o *OrganizationData) HasLocation() bool {
	return o.LocationName != "" || o.Country != "" || o.Region != "" || o.Locality != "" || o.Address != "" || o.Address2 != "" || o.Zip != ""
}

func (o *OrganizationData) HasNotes() bool {
	return len(o.Notes) > 0
}

func (o *OrganizationData) HasPhoneNumbers() bool {
	return len(o.PhoneNumbers) > 0
}

func (o *OrganizationData) HasEmail() bool {
	return o.Email != ""
}

func (o *OrganizationData) HasOwner() bool {
	return o.OwnerUser != nil
}

func (o *OrganizationData) IsSubOrg() bool {
	return o.ParentOrganization != nil && o.ParentOrganization.Organization.Available()
}

func (o *OrganizationData) Normalize() {
	o.SetTimes()
	o.BaseData.Normalize()

	o.NormalizeDomains()

	o.PhoneNumbers = GetNonEmptyPhoneNumbers(o.PhoneNumbers)
	o.PhoneNumbers = RemoveDuplicatedPhoneNumbers(o.PhoneNumbers)

	if strings.TrimSpace(o.CustomerOsId) != "" {
		o.UpdateOnly = true
	}
}

func (o *OrganizationData) NormalizeDomains() {
	o.Domains = utils.FilterOutEmpty(o.Domains)
	o.Domains = utils.LowercaseSliceOfStrings(o.Domains)
	o.Domains = utils.RemoveDuplicates(o.Domains)
}
