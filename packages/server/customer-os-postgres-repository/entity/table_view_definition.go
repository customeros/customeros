package entity

import "time"

type ColumnView struct {
	ColumnId      int    `json:"columnId"`
	ColumnType    string `json:"columnType"`
	Width         int    `json:"width"`
	Visible       bool   `json:"visible"`
	Name          string `json:"name"`
	Filter        string `json:"filter"`
	DefaultFilter string `json:"defaultFilter"`
}

type Columns struct {
	Columns []ColumnView `json:"columns"`
}

type TableViewDefinition struct {
	ID             uint64    `gorm:"primary_key;autoIncrement:true" json:"id"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Tenant         string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	UserId         string    `gorm:"column:user_id;type:varchar(255)" json:"userId"`
	TableId        string    `gorm:"column:table_id;type:varchar(255);NOT NULL;DEFAULT:''" json:"tableId"`
	TableType      string    `gorm:"column:table_type;type:varchar(255);NOT NULL" json:"tableType"`
	Name           string    `gorm:"column:table_name;type:varchar(255);NOT NULL" json:"tableName"`
	Order          int       `gorm:"column:position;type:int;NOT NULL" json:"order"`
	Icon           string    `gorm:"column:icon;type:varchar(255)" json:"icon"`
	Filters        string    `gorm:"column:filters;type:text" json:"filters"`
	DefaultFilters string    `gorm:"column:default_filters;type:text" json:"defaultFilters"`
	Sorting        string    `gorm:"column:sorting;type:text" json:"sorting"`
	ColumnsJson    string    `gorm:"column:columns;type:text" json:"columns"`
	IsPreset       bool      `gorm:"column:is_preset;type:boolean;NOT NULL;DEFAULT:false" json:"isPreset"`
	IsShared       bool      `gorm:"column:is_shared;type:boolean;NOT NULL;DEFAULT:false" json:"isShared"`
}

func (TableViewDefinition) TableName() string {
	return "table_view_definition"
}

type TableViewType string

const (
	TableViewTypeOrganizations TableViewType = "ORGANIZATIONS"
	TableViewTypeInvoices      TableViewType = "INVOICES"
	TableViewTypeContacts      TableViewType = "CONTACTS"
	TableViewTypeOpportunities TableViewType = "OPPORTUNITIES"
	TableViewTypeContracts     TableViewType = "CONTRACTS"
	TableViewTypeFlow          TableViewType = "FLOW"
)

type TableIdType string

const (
	TableIDTypeOrganizations                  TableIdType = "ORGANIZATIONS"
	TableIDTypeCustomers                      TableIdType = "CUSTOMERS"
	TableIDTypeTargets                        TableIdType = "TARGETS"
	TableIDTypeUpcomingInvoices               TableIdType = "UPCOMING_INVOICES"
	TableIDTypePastInvoices                   TableIdType = "PAST_INVOICES"
	TableIDTypeContacts                       TableIdType = "CONTACTS"
	TableIDTypeContactsForTargetOrganizations TableIdType = "CONTACTS_FOR_TARGET_ORGANIZATIONS"
	TableIDTypeOpportunities                  TableIdType = "OPPORTUNITIES"
	TableIDTypeOpportunitiesRecords           TableIdType = "OPPORTUNITIES_RECORDS"
	TableIDTypeContracts                      TableIdType = "CONTRACTS"
	TableIDTypeFlowActions                    TableIdType = "FLOW_ACTIONS"
	TableIDTypeFlowContacts                   TableIdType = "FLOW_CONTACTS"
)

type ColumnViewType string

const (
	ColumnViewTypeInvoicesIssueDate                  ColumnViewType = "INVOICES_ISSUE_DATE"
	ColumnViewTypeInvoicesIssueDatePast              ColumnViewType = "INVOICES_ISSUE_DATE_PAST"
	ColumnViewTypeInvoicesDueDate                    ColumnViewType = "INVOICES_DUE_DATE"
	ColumnViewTypeInvoicesContract                   ColumnViewType = "INVOICES_CONTRACT"
	ColumnViewTypeInvoicesBillingCycle               ColumnViewType = "INVOICES_BILLING_CYCLE"
	ColumnViewTypeInvoicesInvoiceNumber              ColumnViewType = "INVOICES_INVOICE_NUMBER"
	ColumnViewTypeInvoicesAmount                     ColumnViewType = "INVOICES_AMOUNT"
	ColumnViewTypeInvoicesInvoiceStatus              ColumnViewType = "INVOICES_INVOICE_STATUS"
	ColumnViewTypeInvoicesInvoicePreview             ColumnViewType = "INVOICES_INVOICE_PREVIEW"
	ColumnViewTypeInvoicesOrganization               ColumnViewType = "INVOICES_ORGANIZATION"
	ColumnViewTypeOrganizationsAvatar                ColumnViewType = "ORGANIZATIONS_AVATAR"
	ColumnViewTypeOrganizationsName                  ColumnViewType = "ORGANIZATIONS_NAME"
	ColumnViewTypeOrganizationsWebsite               ColumnViewType = "ORGANIZATIONS_WEBSITE"
	ColumnViewTypeOrganizationsPrimaryDomains        ColumnViewType = "ORGANIZATIONS_PRIMARY_DOMAINS"
	ColumnViewTypeOrganizationsRelationship          ColumnViewType = "ORGANIZATIONS_RELATIONSHIP"
	ColumnViewTypeOrganizationsOnboardingStatus      ColumnViewType = "ORGANIZATIONS_ONBOARDING_STATUS"
	ColumnViewTypeOrganizationsRenewalLikelihood     ColumnViewType = "ORGANIZATIONS_RENEWAL_LIKELIHOOD"
	ColumnViewTypeOrganizationsRenewalDate           ColumnViewType = "ORGANIZATIONS_RENEWAL_DATE"
	ColumnViewTypeOrganizationsForecastArr           ColumnViewType = "ORGANIZATIONS_FORECAST_ARR"
	ColumnViewTypeOrganizationsOwner                 ColumnViewType = "ORGANIZATIONS_OWNER"
	ColumnViewTypeOrganizationsLastTouchpoint        ColumnViewType = "ORGANIZATIONS_LAST_TOUCHPOINT"
	ColumnViewTypeOrganizationsLastTouchpointDate    ColumnViewType = "ORGANIZATIONS_LAST_TOUCHPOINT_DATE"
	ColumnViewTypeOrganizationsStage                 ColumnViewType = "ORGANIZATIONS_STAGE"
	ColumnViewTypeOrganizationsContactCount          ColumnViewType = "ORGANIZATIONS_CONTACT_COUNT"
	ColumnViewTypeOrganizationsSocials               ColumnViewType = "ORGANIZATIONS_SOCIALS"
	ColumnViewTypeOrganizationsLeadSource            ColumnViewType = "ORGANIZATIONS_LEAD_SOURCE"
	ColumnViewTypeOrganizationsCreatedDate           ColumnViewType = "ORGANIZATIONS_CREATED_DATE"
	ColumnViewTypeOrganizationsEmployeeCount         ColumnViewType = "ORGANIZATIONS_EMPLOYEE_COUNT"
	ColumnViewTypeOrganizationsYearFounded           ColumnViewType = "ORGANIZATIONS_YEAR_FOUNDED"
	ColumnViewTypeOrganizationsIndustry              ColumnViewType = "ORGANIZATIONS_INDUSTRY"
	ColumnViewTypeOrganizationsChurnDate             ColumnViewType = "ORGANIZATIONS_CHURN_DATE"
	ColumnViewTypeOrganizationsLtv                   ColumnViewType = "ORGANIZATIONS_LTV"
	ColumnViewTypeOrganizationsCountry               ColumnViewType = "ORGANIZATIONS_COUNTRY"
	ColumnViewTypeOrganizationsCity                  ColumnViewType = "ORGANIZATIONS_CITY"
	ColumnViewTypeOrganizationsHeadquarters          ColumnViewType = "ORGANIZATIONS_HEADQUARTERS"
	ColumnViewTypeOrganizationsIsPublic              ColumnViewType = "ORGANIZATIONS_IS_PUBLIC"
	ColumnViewTypeOrganizationsLinkedinFollowerCount ColumnViewType = "ORGANIZATIONS_LINKEDIN_FOLLOWER_COUNT"
	ColumnViewTypeOrganizationsTags                  ColumnViewType = "ORGANIZATIONS_TAGS"
	ColumnViewTypeOrganizationsParentOrganization    ColumnViewType = "ORGANIZATIONS_PARENT_ORGANIZATION"
	ColumnViewTypeOrganizationsUpdatedDate           ColumnViewType = "ORGANIZATIONS_UPDATED_DATE"
	ColumnViewTypeContactsAvatar                     ColumnViewType = "CONTACTS_AVATAR"
	ColumnViewTypeContactsName                       ColumnViewType = "CONTACTS_NAME"
	ColumnViewTypeContactsOrganization               ColumnViewType = "CONTACTS_ORGANIZATION"
	ColumnViewTypeContactsEmails                     ColumnViewType = "CONTACTS_EMAILS"
	ColumnViewTypeContactsPersonalEmails             ColumnViewType = "CONTACTS_PERSONAL_EMAILS"
	ColumnViewTypeContactsPrimaryEmail               ColumnViewType = "CONTACTS_PRIMARY_EMAIL"
	ColumnViewTypeEmailVerificationPrimaryEmail      ColumnViewType = "EMAIL_VERIFICATION_PRIMARY_EMAIL"
	ColumnViewTypeContactsPhoneNumbers               ColumnViewType = "CONTACTS_PHONE_NUMBERS"
	ColumnViewTypeContactsLinkedin                   ColumnViewType = "CONTACTS_LINKEDIN"
	ColumnViewTypeContactsCity                       ColumnViewType = "CONTACTS_CITY"
	ColumnViewTypeContactsPersona                    ColumnViewType = "CONTACTS_PERSONA"
	ColumnViewTypeContactsLastInteraction            ColumnViewType = "CONTACTS_LAST_INTERACTION"
	ColumnViewTypeContactsCountry                    ColumnViewType = "CONTACTS_COUNTRY"
	ColumnViewTypeContactsRegion                     ColumnViewType = "CONTACTS_REGION"
	ColumnViewTypeContactsSkills                     ColumnViewType = "CONTACTS_SKILLS"
	ColumnViewTypeContactsSchools                    ColumnViewType = "CONTACTS_SCHOOLS"
	ColumnViewTypeContactsLanguages                  ColumnViewType = "CONTACTS_LANGUAGES"
	ColumnViewTypeContactsTimeInCurrentRole          ColumnViewType = "CONTACTS_TIME_IN_CURRENT_ROLE"
	ColumnViewTypeContactsExperience                 ColumnViewType = "CONTACTS_EXPERIENCE"
	ColumnViewTypeContactsLinkedinFollowerCount      ColumnViewType = "CONTACTS_LINKEDIN_FOLLOWER_COUNT"
	ColumnViewTypeContactsJobTitle                   ColumnViewType = "CONTACTS_JOB_TITLE"
	ColumnViewTypeContactsTags                       ColumnViewType = "CONTACTS_TAGS"
	ColumnViewTypeContactsConnections                ColumnViewType = "CONTACTS_CONNECTIONS"
	ColumnViewTypeContactsFlows                      ColumnViewType = "CONTACTS_FLOWS"
	ColumnViewTypeContactsFlowStatus                 ColumnViewType = "CONTACTS_FLOW_STATUS"
	ColumnViewTypeContactsFlowNextAction             ColumnViewType = "CONTACTS_FLOW_NEXT_ACTION"
	ColumnViewTypeContactsUpdatedAt                  ColumnViewType = "CONTACTS_UPDATED_AT"
	ColumnViewTypeContactsCreatedAt                  ColumnViewType = "CONTACTS_CREATED_AT"
	ColumnViewTypeOpportunitiesCommonColumn          ColumnViewType = "OPPORTUNITIES_COMMON_COLUMN"
	ColumnViewTypeOpportunitiesName                  ColumnViewType = "OPPORTUNITIES_NAME"
	ColumnViewTypeOpportunitiesOrganization          ColumnViewType = "OPPORTUNITIES_ORGANIZATION"
	ColumnViewTypeOpportunitiesStage                 ColumnViewType = "OPPORTUNITIES_STAGE"
	ColumnViewTypeOpportunitiesEstimatedArr          ColumnViewType = "OPPORTUNITIES_ESTIMATED_ARR"
	ColumnViewTypeOpportunitiesOwner                 ColumnViewType = "OPPORTUNITIES_OWNER"
	ColumnViewTypeOpportunitiesTimeInStage           ColumnViewType = "OPPORTUNITIES_TIME_IN_STAGE"
	ColumnViewTypeOpportunitiesCreatedDate           ColumnViewType = "OPPORTUNITIES_CREATED_DATE"
	ColumnViewTypeOpportunitiesNextStep              ColumnViewType = "OPPORTUNITIES_NEXT_STEP"
	ColumnViewTypeContractsName                      ColumnViewType = "CONTRACTS_NAME"
	ColumnViewTypeContractsEnded                     ColumnViewType = "CONTRACTS_ENDED"
	ColumnViewTypeContractsPeriod                    ColumnViewType = "CONTRACTS_PERIOD"
	ColumnViewTypeContractsCurrency                  ColumnViewType = "CONTRACTS_CURRENCY"
	ColumnViewTypeContractsStatus                    ColumnViewType = "CONTRACTS_STATUS"
	ColumnViewTypeContractsRenewal                   ColumnViewType = "CONTRACTS_RENEWAL"
	ColumnViewTypeContractsLtv                       ColumnViewType = "CONTRACTS_LTV"
	ColumnViewTypeContractsRenewalDate               ColumnViewType = "CONTRACTS_RENEWAL_DATE"
	ColumnViewTypeContractsForecastArr               ColumnViewType = "CONTRACTS_FORECAST_ARR"
	ColumnViewTypeContractsOwner                     ColumnViewType = "CONTRACTS_OWNER"
	ColumnViewTypeContractsHealth                    ColumnViewType = "CONTRACTS_HEALTH"
	ColumnViewTypeFlowName                           ColumnViewType = "FLOW_NAME"
	ColumnViewTypeFlowTotalCount                     ColumnViewType = "FLOW_TOTAL_COUNT"
	ColumnViewTypeFlowOnHoldCount                    ColumnViewType = "FLOW_ON_HOLD_COUNT"
	ColumnViewTypeFlowReadyCount                     ColumnViewType = "FLOW_READY_COUNT"
	ColumnViewTypeFlowScheduledCount                 ColumnViewType = "FLOW_SCHEDULED_COUNT"
	ColumnViewTypeFlowInProgressCount                ColumnViewType = "FLOW_IN_PROGRESS_COUNT"
	ColumnViewTypeFlowCompletedCount                 ColumnViewType = "FLOW_COMPLETED_COUNT"
	ColumnViewTypeFlowGoalAchievedCount              ColumnViewType = "FLOW_GOAL_ACHIEVED_COUNT"
	ColumnViewTypeFlowStatus                         ColumnViewType = "FLOW_STATUS"
	ColumnViewTypeFlowActionName                     ColumnViewType = "FLOW_ACTION_NAME"
	ColumnViewTypeFlowActionStatus                   ColumnViewType = "FLOW_ACTION_STATUS"
)
