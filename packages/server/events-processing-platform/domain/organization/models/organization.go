package models

import (
	"fmt"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type RenewalLikelihoodProbability string

const (
	RenewalLikelihoodHIGH   RenewalLikelihoodProbability = "HIGH"
	RenewalLikelihoodMEDIUM RenewalLikelihoodProbability = "MEDIUM"
	RenewalLikelihoodLOW    RenewalLikelihoodProbability = "LOW"
	RenewalLikelihoodZERO   RenewalLikelihoodProbability = "ZERO"
)

type CustomFieldDataType string

const (
	CustomFieldDataTypeText     CustomFieldDataType = "TEXT"
	CustomFieldDataTypeBool     CustomFieldDataType = "BOOL"
	CustomFieldDataTypeDatetime CustomFieldDataType = "DATETIME"
	CustomFieldDataTypeInteger  CustomFieldDataType = "INTEGER"
	CustomFieldDataTypeDecimal  CustomFieldDataType = "DECIMAL"
)

func (r RenewalLikelihoodProbability) CamelCaseString() string {
	switch r {
	case RenewalLikelihoodHIGH:
		return "High"
	case RenewalLikelihoodMEDIUM:
		return "Medium"
	case RenewalLikelihoodLOW:
		return "Low"
	case RenewalLikelihoodZERO:
		return "Zero"
	}
	return ""
}

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
	RenewalLikelihood   RenewalLikelihood                  `json:"renewalLikelihood,omitempty"`
	RenewalForecast     RenewalForecast                    `json:"renewalForecast,omitempty"`
	BillingDetails      BillingDetails                     `json:"billingDetails,omitempty"`
	CustomFields        map[string]CustomField             `json:"customFields,omitempty"`
	ExternalSystems     []cmnmod.ExternalSystem            `json:"externalSystems"`
	ParentOrganizations map[string]ParentOrganization      `json:"parentOrganizations,omitempty"`
}

type RenewalLikelihood struct {
	RenewalLikelihood RenewalLikelihoodProbability `json:"renewalLikelihood,omitempty"`
	Comment           *string                      `json:"comment,omitempty"`
	UpdatedAt         time.Time                    `json:"updatedAt,omitempty"`
	UpdatedBy         string                       `json:"updatedBy,omitempty"`
}

type RenewalForecast struct {
	Amount    *float64  `json:"amount,omitempty"`
	Comment   *string   `json:"comment,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	UpdatedBy string    `json:"updatedBy,omitempty"`
}

type BillingDetails struct {
	Amount            *float64   `json:"amount,omitempty"`
	UpdatedBy         string     `json:"updatedBy,omitempty"`
	Frequency         string     `json:"frequency,omitempty"`
	RenewalCycle      string     `json:"renewalCycle,omitempty"`
	RenewalCycleStart *time.Time `json:"renewalCycleStart,omitempty"`
	RenewalCycleNext  *time.Time `json:"renewalCycleNext,omitempty"`
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
