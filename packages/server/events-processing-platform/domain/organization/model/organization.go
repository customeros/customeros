package model

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"reflect"
	"time"

	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
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
	FieldMaskEmployees          = "employees"
	FieldMaskLastFundingRound   = "lastFundingRound"
	FieldMaskLastFundingAmount  = "lastFundingAmount"
	FieldMaskReferenceId        = "referenceId"
	FieldMaskWebsite            = "website"
	FieldMaskYearFounded        = "yearFounded"
	FieldMaskHeadquarters       = "headquarters"
	FieldMaskLogoUrl            = "logoUrl"
	FieldMaskIconUrl            = "iconUrl"
	FieldMaskEmployeeGrowthRate = "employeeGrowthRate"
	FieldMaskSlackChannelId     = "slackChannelId"
	FieldMaskRelationship       = "relationship"
	FieldMaskStage              = "stage"
	FieldMaskIcpFit             = "icpFit"
)

type CustomFieldDataType string

const (
	CustomFieldDataTypeText     CustomFieldDataType = "TEXT"
	CustomFieldDataTypeBool     CustomFieldDataType = "BOOL"
	CustomFieldDataTypeDatetime CustomFieldDataType = "DATETIME"
	CustomFieldDataTypeInteger  CustomFieldDataType = "INTEGER"
	CustomFieldDataTypeDecimal  CustomFieldDataType = "DECIMAL"
)

type CustomField struct {
	Id                  string                      `json:"id"`
	Name                string                      `json:"name"`
	TemplateId          *string                     `json:"templateId,omitempty"`
	CustomFieldValue    neo4jmodel.CustomFieldValue `json:"customFieldValue"`
	CustomFieldDataType CustomFieldDataType         `json:"customFieldDataType"`
	Source              common.Source               `json:"source"`
	CreatedAt           time.Time                   `json:"createdAt,omitempty"`
	UpdatedAt           time.Time                   `json:"updatedAt,omitempty"`
}

type Organization struct {
	ID                string                             `json:"id"`
	Name              string                             `json:"name"`
	Hide              bool                               `json:"hide"`
	Description       string                             `json:"description"`
	Website           string                             `json:"website"`
	Industry          string                             `json:"industry"`
	SubIndustry       string                             `json:"subIndustry"`
	IndustryGroup     string                             `json:"industryGroup"`
	TargetAudience    string                             `json:"targetAudience"`
	ValueProposition  string                             `json:"valueProposition"`
	IsPublic          bool                               `json:"isPublic"`
	Employees         int64                              `json:"employees"`
	Market            string                             `json:"market"`
	LastFundingRound  string                             `json:"lastFundingRound"`
	LastFundingAmount string                             `json:"lastFundingAmount"`
	ReferenceId       string                             `json:"referenceId"`
	Note              string                             `json:"note"`
	Source            common.Source                      `json:"source"`
	CreatedAt         time.Time                          `json:"createdAt,omitempty"`
	UpdatedAt         time.Time                          `json:"updatedAt,omitempty"`
	PhoneNumbers      map[string]OrganizationPhoneNumber `json:"phoneNumbers"`
	// Deprecated
	LocationIds         []string                      `json:"locationIds,omitempty"`
	Socials             map[string]common.Social      `json:"socials,omitempty"`
	CustomFields        map[string]CustomField        `json:"customFields,omitempty"`
	ExternalSystems     []common.ExternalSystem       `json:"externalSystems"`
	ParentOrganizations map[string]ParentOrganization `json:"parentOrganizations,omitempty"`
	LogoUrl             string                        `json:"logoUrl,omitempty"`
	IconUrl             string                        `json:"iconUrl,omitempty"`
	YearFounded         *int64                        `json:"yearFounded,omitempty"`
	Headquarters        string                        `json:"headquarters,omitempty"`
	EmployeeGrowthRate  string                        `json:"employeeGrowthRate,omitempty"`
	SlackChannelId      string                        `json:"slackChannelId,omitempty"`
	OnboardingDetails   OnboardingDetails             `json:"onboardingDetails,omitempty"`
	BillingProfiles     map[string]BillingProfile     `json:"billingProfiles,omitempty"`
	Relationship        string                        `json:"relationship,omitempty"`
	Stage               string                        `json:"stage,omitempty"`
	LeadSource          string                        `json:"leadSource,omitempty"`
	TagIds              []string                      `json:"tagIds,omitempty"`
	Locations           map[string]common.Location    `json:"locations,omitempty"`
	IcpFit              bool                          `json:"icpFit"`
}

type BillingProfile struct {
	Id             string        `json:"id"`
	LegalName      string        `json:"legalName"`
	TaxId          string        `json:"taxId"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	SourceFields   common.Source `json:"sourceFields"`
	PrimaryEmailId string        `json:"primaryEmailId"`
	EmailIds       []string      `json:"emailIds"`
	LocationIds    []string      `json:"locationIds"`
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

type ParentOrganization struct {
	OrganizationId string `json:"organizationId"`
	Type           string `json:"type"`
}

func (o *Organization) String() string {
	return fmt.Sprintf("Organization{ID: %s, Name: %s, Description: %s, Website: %s, Industry: %s, IsPublic: %t, SourceFields: %s, CreatedAt: %s, UpdatedAt: %s}", o.ID, o.Name, o.Description, o.Website, o.Industry, o.IsPublic, o.Source, o.CreatedAt, o.UpdatedAt)
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

func (o *Organization) GetLocationIdForDetails(location common.Location) string {
	for id, orgLocation := range o.Locations {
		if locationMatchesExcludingName(orgLocation, location) {
			return id
		}
	}
	return ""
}

func locationMatchesExcludingName(orgLocation, inputLocation common.Location) bool {
	// Create copies of the locations to avoid modifying the original structs
	orgCopy := orgLocation
	inputCopy := inputLocation

	// Set Name to empty string for both locations to exclude it from comparison
	orgCopy.Name = ""
	inputCopy.Name = ""

	// Compare all fields except Name
	return reflect.DeepEqual(orgCopy, inputCopy)
}

func (o *Organization) ContainsExternalSystem(externalSystem string) bool {
	for _, es := range o.ExternalSystems {
		if es.ExternalSystemId == externalSystem {
			return true
		}
	}
	return false
}

func (o *Organization) SkipUpdate(fields *OrganizationFields) bool {
	if fields.ExternalSystem.Available() && !o.ContainsExternalSystem(fields.ExternalSystem.ExternalSystemId) {
		return false
	}
	if o.Source.SourceOfTruth == constants.SourceOpenline && fields.ExternalSystem.Available() {
		return true
	}
	return false
}
