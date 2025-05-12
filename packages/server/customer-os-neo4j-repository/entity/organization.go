package neo4j_entity

import (
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type OrganizationProperty string

const (
	OrganizationPropertyEmployees                 OrganizationProperty = "employees"
	OrganizationPropertyUpdatedAt                 OrganizationProperty = "updatedAt"
	OrganizationPropertyYearFounded               OrganizationProperty = "yearFounded"
	OrganizationPropertyName                      OrganizationProperty = "name"
	OrganizationPropertyHide                      OrganizationProperty = "hide"
	OrganizationPropertyStage                     OrganizationProperty = "stage"
	OrganizationPropertyIndustry                  OrganizationProperty = "industry"
	OrganizationPropertyWrongIndustry             OrganizationProperty = "wrongIndustry"
	OrganizationPropertyWrongIndustryAt           OrganizationProperty = "wrongIndustryAt"
	OrganizationPropertyIsPublic                  OrganizationProperty = "isPublic"
	OrganizationPropertyDomainCheckedAt           OrganizationProperty = "techDomainCheckedAt"
	OrganizationPropertyLastTouchpointRequestedAt OrganizationProperty = "techLastTouchpointRequestedAt"
	OrganizationPropertyLastTouchpointType        OrganizationProperty = "lastTouchpointType"
	OrganizationPropertyLastTouchpointAt          OrganizationProperty = "lastTouchpointAt"
	OrganizationPropertyIcpFit                    OrganizationProperty = "icpFit"
	OrganizationPropertyIcpFitUpdatedAt           OrganizationProperty = "icpFitUpdatedAt"
	OrganizationPropertyIcpFitReasons             OrganizationProperty = "icpFitReasons"
	OrganizationPropertyHiddenAt                  OrganizationProperty = "hiddenAt"
	OrganizationPropertyEnrichRequestedAt         OrganizationProperty = "techEnrichRequestedAt"
	OrganizationPropertyEnrichedAt                OrganizationProperty = "enrichedAt"
	OrganizationPropertyEnrichFailedAt            OrganizationProperty = "enrichFailedAt"
	OrganizationPropertyEnrichAttempts            OrganizationProperty = "techEnrichAttempts"
	OrganizationPropertyRenewalLikelihood         OrganizationProperty = "derivedRenewalLikelihood"
	OrganizationPropertyContactCount              OrganizationProperty = "derivedContactCount"
	OrganizationPropertyLeadSource                OrganizationProperty = "leadSource"
	OrganizationPropertyIcpCheckedAt              OrganizationProperty = "techIcpCheckedAt"
	OrganizationPropertyIcpCheckRequestedAt       OrganizationProperty = "techIcpCheckRequestedAt"
)

type OrganizationEntity struct {
	EventStoreAggregate // TODO alexb delete
	DataLoaderKey
	ID                 string
	CustomerOsId       string `neo4jDb:"property:customerOsId;lookupName:CUSTOMER_OS_ID;supportCaseSensitive:false"`
	CustomId           string `neo4jDb:"property:customId;lookupName:CUSTOMER_ID;supportCaseSensitive:false"`
	Name               string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	Description        string `neo4jDb:"property:description;lookupName:DESCRIPTION;supportCaseSensitive:true"`
	Website            string `neo4jDb:"property:website;lookupName:WEBSITE;supportCaseSensitive:true"`
	Industry           string // Free text industry, replaced with link to industry node. Use industry as temporary field for sync from other systems
	IsPublic           bool
	Hide               bool
	Market             string
	LastFundingRound   string
	LastFundingAmount  string
	ReferenceId        string `neo4jDb:"property:referenceId;lookupName:REFERENCE_ID;supportCaseSensitive:true"`
	Note               string
	Employees          int64
	CreatedAt          time.Time
	UpdatedAt          time.Time  `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	LastTouchpointAt   *time.Time `neo4jDb:"property:lastTouchpointAt;lookupName:LAST_TOUCHPOINT_AT;supportCaseSensitive:false"`
	LastTouchpointId   *string    `neo4jDb:"property:lastTouchpointId;lookupName:LAST_TOUCHPOINT_ID;supportCaseSensitive:false"`
	LastTouchpointType *string    `neo4jDb:"property:lastTouchpointType;lookupName:LAST_TOUCHPOINT_TYPE;supportCaseSensitive:false"`
	Source             DataSource
	YearFounded        *int64
	Headquarters       string
	EmployeeGrowthRate string
	SlackChannelId     string
	LogoUrl            string
	IconUrl            string
	Stage              commonenum.OrganizationStage `neo4jDb:"property:stage;lookupName:STAGE;supportCaseSensitive:false"`
	StageUpdatedAt     *time.Time
	LeadSource         string `neo4jDb:"property:leadSource;lookupName:LEAD_SOURCE;supportCaseSensitive:true"`
	IcpFit             commonenum.IcpFit
	IcpFitUpdatedAt    *time.Time
	IcpFitReasons      []string
	WrongIndustry      bool

	LinkedOrganizationType *string

	SuggestedMerge struct {
		SuggestedAt *time.Time
		SuggestedBy *string
		Confidence  *float64
	}
	RenewalSummary                     RenewalSummary
	OnboardingDetails                  OnboardingDetails
	EnrichDetails                      OrganizationEnrichDetails
	InteractionEventParticipantDetails InteractionEventParticipantDetails
	OrganizationInternalFields         OrganizationInternalFields
	DerivedData                        DerivedData

	// Deprecated
	Relationship enum.OrganizationRelationship `neo4jDb:"property:relationship;lookupName:RELATIONSHIP;supportCaseSensitive:false"`
}

type DerivedData struct {
	ChurnedAt    *time.Time    `neo4jDb:"property:derivedNextRenewalAt;lookupName:CHURN_DATE;supportCaseSensitive:false"`
	Ltv          float64       `neo4jDb:"property:derivedLtv;lookupName:LTV;supportCaseSensitive:false"`
	LtvCurrency  enum.Currency `neo4jDb:"property:derivedLtvCurrency;lookupName:LTV_CURRENCY;supportCaseSensitive:false"`
	ContactCount int64         `neo4jDb:"property:derivedContactCount;lookupName:CONTACT_COUNT;supportCaseSensitive:false"`
}

type RenewalSummary struct {
	ArrForecast            *float64
	MaxArrForecast         *float64
	NextRenewalAt          *time.Time `neo4jDb:"property:derivedNextRenewalAt;lookupName:RENEWAL_DATE;supportCaseSensitive:false"`
	RenewalLikelihood      string     `neo4jDb:"property:derivedRenewalLikelihood;lookupName:RENEWAL_LIKELIHOOD;supportCaseSensitive:false"`
	RenewalLikelihoodOrder *int64
}

type OnboardingDetails struct {
	Status       string
	SortingOrder *int64
	UpdatedAt    *time.Time
	Comments     string
}

type OrganizationEnrichDetails struct {
	EnrichRequestedAt *time.Time
	EnrichedAt        *time.Time
	EnrichFailedAt    *time.Time
	EnrichAttempts    int64
	EnrichDomain      string
	EnrichSource      enum.EnrichSource
}

type OrganizationInternalFields struct {
	DomainCheckedAt           *time.Time
	LastTouchpointRequestedAt *time.Time
	HiddenAt                  *time.Time
	IcpCheckedAt              *time.Time
	IcpCheckRequestedAt       *time.Time
}

type OrganizationEntities []OrganizationEntity

func (OrganizationEntity) IsInteractionEventParticipant() {}

func (OrganizationEntity) IsIssueParticipant() {}

func (OrganizationEntity) IsNotedEntity() {}

func (OrganizationEntity) IsMeetingParticipant() {}

func (OrganizationEntity) EntityLabel() string {
	return model.NodeLabelOrganization
}

func (o OrganizationEntity) GetDataloaderKey() string {
	return o.DataloaderKey
}

func (o OrganizationEntity) Labels(tenant string) []string {
	return []string{
		o.EntityLabel(),
		o.EntityLabel() + "_" + tenant,
	}
}

func (o OrganizationEntity) IsHidden() bool {
	return o.Hide
}
