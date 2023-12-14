package model

import (
	"fmt"
	"time"

	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
)

const (
	FieldMaskName               = "name"
	FieldMaskTargetAudience     = "targetAudience"
	FieldMaskValueProposition   = "valueProposition"
	FieldMaskIndustry           = "industry"
	FieldMaskSubIndustry        = "subIndustry"
	FieldMaskIndustryGroup      = "industryGroup"
	FieldMaskMarket             = "market"
	FieldMaskHide               = "hide"
	FieldMaskDescription        = "description"
	FieldMaskNote               = "note"
	FieldMaskIsPublic           = "isPublic"
	FieldMaskIsCustomer         = "isCustomer"
	FieldMaskEmployees          = "employees"
	FieldMaskLastFundingRound   = "lastFundingRound"
	FieldMaskLastFundingAmount  = "lastFundingAmount"
	FieldMaskReferenceId        = "referenceId"
	FieldMaskWebsite            = "website"
	FieldMaskYearFounded        = "yearFounded"
	FieldMaskHeadquarters       = "headquarters"
	FieldMaskLogoUrl            = "logoUrl"
	FieldMaskEmployeeGrowthRate = "employeeGrowthRate"
)

type CustomFieldDataType string

const (
	CustomFieldDataTypeText     CustomFieldDataType = "TEXT"
	CustomFieldDataTypeBool     CustomFieldDataType = "BOOL"
	CustomFieldDataTypeDatetime CustomFieldDataType = "DATETIME"
	CustomFieldDataTypeInteger  CustomFieldDataType = "INTEGER"
	CustomFieldDataTypeDecimal  CustomFieldDataType = "DECIMAL"
)

type Social struct {
	PlatformName string `json:"platformName"`
	Url          string `json:"url"`
}

type CustomFieldValue struct {
	Str     *string    `json:"string,omitempty"`
	Int     *int64     `json:"int,omitempty"`
	Time    *time.Time `json:"time,omitempty"`
	Bool    *bool      `json:"bool,omitempty"`
	Decimal *float64   `json:"decimal,omitempty"`
}

func (c *CustomFieldValue) RealValue() any {
	if c.Int != nil {
		return *c.Int
	} else if c.Decimal != nil {
		return *c.Decimal
	} else if c.Time != nil {
		return *c.Time
	} else if c.Bool != nil {
		return *c.Bool
	} else if c.Str != nil {
		return *c.Str
	}
	return nil
}

type CustomField struct {
	Id                  string              `json:"id"`
	Name                string              `json:"name"`
	TemplateId          *string             `json:"templateId,omitempty"`
	CustomFieldValue    CustomFieldValue    `json:"customFieldValue"`
	CustomFieldDataType CustomFieldDataType `json:"customFieldDataType"`
	Source              cmnmod.Source       `json:"source"`
	CreatedAt           time.Time           `json:"createdAt,omitempty"`
	UpdatedAt           time.Time           `json:"updatedAt,omitempty"`
}

type Organization struct {
	ID                  string                             `json:"id"`
	Name                string                             `json:"name"`
	Hide                bool                               `json:"hide"`
	Description         string                             `json:"description"`
	Website             string                             `json:"website"`
	Industry            string                             `json:"industry"`
	SubIndustry         string                             `json:"subIndustry"`
	IndustryGroup       string                             `json:"industryGroup"`
	TargetAudience      string                             `json:"targetAudience"`
	ValueProposition    string                             `json:"valueProposition"`
	IsPublic            bool                               `json:"isPublic"`
	IsCustomer          bool                               `json:"isCustomer"`
	Employees           int64                              `json:"employees"`
	Market              string                             `json:"market"`
	LastFundingRound    string                             `json:"lastFundingRound"`
	LastFundingAmount   string                             `json:"lastFundingAmount"`
	ReferenceId         string                             `json:"referenceId"`
	Note                string                             `json:"note"`
	Source              cmnmod.Source                      `json:"source"`
	CreatedAt           time.Time                          `json:"createdAt,omitempty"`
	UpdatedAt           time.Time                          `json:"updatedAt,omitempty"`
	PhoneNumbers        map[string]OrganizationPhoneNumber `json:"phoneNumbers"`
	Emails              map[string]OrganizationEmail       `json:"emails"`
	Locations           []string                           `json:"locations,omitempty"`
	Domains             []string                           `json:"domains,omitempty"`
	Socials             map[string]Social                  `json:"socials,omitempty"`
	CustomFields        map[string]CustomField             `json:"customFields,omitempty"`
	ExternalSystems     []cmnmod.ExternalSystem            `json:"externalSystems"`
	ParentOrganizations map[string]ParentOrganization      `json:"parentOrganizations,omitempty"`
	LogoUrl             string                             `json:"logoUrl,omitempty"`
	YearFounded         *int64                             `json:"yearFounded,omitempty"`
	Headquarters        string                             `json:"headquarters,omitempty"`
	EmployeeGrowthRate  string                             `json:"employeeGrowthRate,omitempty"`
	OnboardingDetails   OnboardingDetails                  `json:"onboardingDetails,omitempty"`
}

type OnboardingDetails struct {
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
	Comments  string    `json:"comments"`
}

type OrganizationPhoneNumber struct {
	Primary bool   `json:"primary"`
	Label   string `json:"label"`
}

type OrganizationEmail struct {
	Primary bool   `json:"primary"`
	Label   string `json:"label"`
}

type ParentOrganization struct {
	OrganizationId string `json:"organizationId"`
	Type           string `json:"type"`
}

func (o *Organization) String() string {
	return fmt.Sprintf("Organization{ID: %s, Name: %s, Description: %s, Website: %s, Industry: %s, IsPublic: %t, Source: %s, CreatedAt: %s, UpdatedAt: %s}", o.ID, o.Name, o.Description, o.Website, o.Industry, o.IsPublic, o.Source, o.CreatedAt, o.UpdatedAt)
}

func (o *Organization) GetSocialIdForUrl(url string) string {
	if o.Socials == nil {
		return ""
	}
	for key, social := range o.Socials {
		if social.Url == url {
			return key
		}
	}
	return ""
}
