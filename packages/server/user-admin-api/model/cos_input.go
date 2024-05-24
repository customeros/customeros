package model

import (
	"time"
)

type WorkspaceInput struct {
	Name      string  `json:"name"`
	Provider  string  `json:"provider"`
	AppSource *string `json:"appSource"`
}

type EmailInput struct {
	Email     string  `json:"email"`
	Primary   bool    `json:"primary"`
	AppSource *string `json:"appSource"`
}

type PlayerInput struct {
	IdentityId string  `json:"identityId"`
	AuthId     string  `json:"authId"`
	Provider   string  `json:"provider"`
	AppSource  *string `json:"appSource"`
}

type UserInput struct {
	FirstName       string      `json:"firstName"`
	LastName        string      `json:"lastName"`
	Email           EmailInput  `json:"email"`
	Player          PlayerInput `json:"player"`
	AppSource       *string     `json:"appSource"`
	ProfilePhotoURL *string     `json:"profilePhotoUrl,omitempty"`
}

type TenantInput struct {
	Name      string  `json:"name"`
	AppSource *string `json:"appSource"`
}
type TenantBillingProfileInput struct {
	LegalName                     string `json:"legalName"`
	Email                         string `json:"email"`
	AddressLine1                  string `json:"addressLine1"`
	Locality                      string `json:"locality"`
	Country                       string `json:"country"`
	Zip                           string `json:"zip"`
	DomesticPaymentsBankInfo      string `json:"domesticPaymentsBankInfo"`
	InternationalPaymentsBankInfo string `json:"internationalPaymentsBankInfo"`
	VatNumber                     string `json:"vatNumber"`
	SendInvoicesFrom              string `json:"sendInvoicesFrom"`
	CanPayWithCard                bool   `json:"canPayWithCard"`
	CanPayWithDirectDebitSEPA     bool   `json:"canPayWithDirectDebitSEPA"`
	CanPayWithDirectDebitACH      bool   `json:"canPayWithDirectDebitACH"`
	CanPayWithDirectDebitBacs     bool   `json:"canPayWithDirectDebitBacs"`
	CanPayWithPigeon              bool   `json:"canPayWithPigeon"`
	CanPayWithBankTransfer        bool   `json:"canPayWithBankTransfer"`
	Check                         bool   `json:"check"`
}

type NextDryRunInvoiceForContractInput struct {
	ContractId string `json:"contractId"`
}

type ContactInput struct {
	FirstName       *string           `json:"firstName,omitempty"`
	LastName        *string           `json:"lastName,omitempty"`
	Email           *EmailInput       `json:"email,omitempty"`
	ProfilePhotoURL *string           `json:"profilePhotoUrl,omitempty"`
	Timezone        *string           `json:"timezone,omitempty"`
	PhoneNumber     *PhoneNumberInput `json:"phoneNumber,omitempty"`
	Description     *string           `json:"description,omitempty"`
}

type PhoneNumberInput struct {
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type SocialInput struct {
	Url string `json:"url,omitempty"`
}

type ContractInput struct {
	OrganizationId          string    `json:"organizationId"`
	ContractName            string    `json:"contractName"`
	CommittedPeriodInMonths int64     `json:"committedPeriodInMonths"`
	ContractUrl             string    `json:"contractUrl"`
	ServiceStarted          time.Time `json:"serviceStarted"`
	ContractSigned          time.Time `json:"contractSigned"`
}

type ContractUpdateInput struct {
	ContractId            string     `json:"contractId,omitempty"`
	Patch                 bool       `json:"patch"`
	EndedAt               time.Time  `json:"endedAt"`
	Currency              string     `json:"currency"`
	InvoicingStartDate    *time.Time `json:"invoicingStartDate"`
	BillingCycle          string     `json:"billingCycle"`
	AddressLine1          string     `json:"addressLine1"`
	AddressLine2          string     `json:"addressLine2"`
	Locality              string     `json:"locality"`
	Country               string     `json:"country"`
	Zip                   string     `json:"zip"`
	OrganizationLegalName string     `json:"organizationLegalName"`
	InvoiceEmail          string     `json:"invoiceEmail"`
	InvoiceNote           string     `json:"invoiceNote"`
}

type ContactOrganizationInput struct {
	ContactId      string `json:"contactId,omitempty"`
	OrganizationId string `json:"organizationId,omitempty"`
}

type ServiceLineInput struct {
	ContractId     string     `json:"contractId,omitempty"`
	Description    string     `json:"description"`
	BillingCycle   string     `json:"billingCycle"`
	Price          int        `json:"price"`
	Quantity       int        `json:"quantity"`
	ServiceStarted *time.Time `json:"serviceStarted"`
}

type ServiceLineEndedInput struct {
	ContractId     string     `json:"contractId,omitempty"`
	Description    string     `json:"description"`
	BillingCycle   string     `json:"billingCycle"`
	Price          int        `json:"price"`
	Quantity       int        `json:"quantity"`
	ServiceStarted *time.Time `json:"serviceStarted"`
	ServiceEnded   *time.Time `json:"serviceEnded"`
}

type InteractionSessionParticipantInput struct {
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	ContactID   *string `json:"contactID,omitempty"`
	UserID      *string `json:"userID,omitempty"`
	Type        *string `json:"type,omitempty"`
}

type InteractionEventParticipantInput struct {
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phoneNumber,omitempty"`
	ContactID   *string `json:"contactID,omitempty"`
	UserID      *string `json:"userID,omitempty"`
	Type        *string `json:"type,omitempty"`
}

type EmailChannelData struct {
	Subject   string   `json:"Subject"`
	InReplyTo []string `json:"InReplyTo"`
	Reference []string `json:"Reference"`
}

type MeetingParticipantInput struct {
	ContactID      *string `json:"contactId,omitempty"`
	UserID         *string `json:"userId,omitempty"`
	OrganizationID *string `json:"organizationId,omitempty"`
}

type NoteInput struct {
	Content     *string `json:"content,omitempty"`
	ContentType *string `json:"contentType,omitempty"`
	AppSource   *string `json:"appSource,omitempty"`
}

type MeetingStatus string

const (
	MeetingStatusUndefined MeetingStatus = "UNDEFINED"
	MeetingStatusAccepted  MeetingStatus = "ACCEPTED"
	MeetingStatusCanceled  MeetingStatus = "CANCELED"
)

type MeetingInput struct {
	Name               *string                    `json:"name,omitempty"`
	AttendedBy         []*MeetingParticipantInput `json:"attendedBy,omitempty"`
	CreatedBy          []*MeetingParticipantInput `json:"createdBy,omitempty"`
	CreatedAt          *time.Time                 `json:"createdAt,omitempty"`
	StartedAt          *time.Time                 `json:"startedAt,omitempty"`
	EndedAt            *time.Time                 `json:"endedAt,omitempty"`
	ConferenceURL      *string                    `json:"conferenceUrl,omitempty"`
	MeetingExternalURL *string                    `json:"meetingExternalUrl,omitempty"`
	Agenda             *string                    `json:"agenda,omitempty"`
	AgendaContentType  *string                    `json:"agendaContentType,omitempty"`
	Note               *NoteInput                 `json:"note,omitempty"`
	AppSource          *string                    `json:"appSource,omitempty"`
	Status             *MeetingStatus             `json:"status,omitempty"`
}

type OrganizationUpdateOnboardingStatus struct {
	OrganizationId string `json:"organizationId"`
	Status         string `json:"status"`
	Comments       string `json:"comments,omitempty"`
}

type MasterPlanMilestoneInput struct {
	MasterPlanId  string   `json:"masterPlanId"`
	Name          string   `json:"name"`
	Order         int64    `json:"order"`
	DurationHours int64    `json:"durationHours"`
	Optional      bool     `json:"optional"`
	Items         []string `json:"items"`
}

type OrganizationRelationship string

const (
	OrganizationRelationshipCustomer       OrganizationRelationship = "CUSTOMER"
	OrganizationRelationshipProspect       OrganizationRelationship = "PROSPECT"
	OrganizationRelationshipStranger       OrganizationRelationship = "STRANGER"
	OrganizationRelationshipFormerCustomer OrganizationRelationship = "FORMER_CUSTOMER"
)

type OrganizationStage string

const (
	OrganizationStageLead       OrganizationStage = "LEAD"
	OrganizationStageTarget     OrganizationStage = "TARGET"
	OrganizationStageInterested OrganizationStage = "INTERESTED"
	OrganizationStageEngaged    OrganizationStage = "ENGAGED"
	OrganizationStageContracted OrganizationStage = "CONTRACTED"
	OrganizationStageNurture    OrganizationStage = "NURTURE"
	OrganizationStageAbandoned  OrganizationStage = "ABANDONED"
)

type OrganizationInput struct {
	Name         *string                   `json:"name,omitempty"`
	Relationship *OrganizationRelationship `json:"relationship,omitempty"`
	Stage        *OrganizationStage        `json:"stage,omitempty"`
	LeadSource   *string                   `json:"leadSource,omitempty"`
	Domains      []string                  `json:"domains,omitempty"`
	Notes        *string                   `json:"notes,omitempty"`
	Industry     *string                   `json:"industry,omitempty"`
	Market       *string                   `json:"market,omitempty"`
	Employees    *int64                    `json:"employees,omitempty"`
	Website      *string                   `json:"website,omitempty"`
}

type OrganizationUpdateInput struct {
	Id                string  `json:"id,omitempty"`
	LastFundingAmount *string `json:"lastFundingAmount,omitempty"`
	LastFundingRound  *string `json:"lastFundingRound,omitempty"`
	TargetAudience    *string `json:"targetAudience,omitempty"`
	ValueProposition  *string `json:"valueProposition,omitempty"`
}
