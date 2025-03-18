package neo4j_entity

import (
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"time"
)

type InvoiceProperty string

const (
	InvoicePropertyCreatedAt                            InvoiceProperty = "createdAt"
	InvoicePropertyFinalizedWebhookProcessedAt          InvoiceProperty = "techInvoiceFinalizedWebhookProcessedAt"
	InvoicePropertyPaidWebhookProcessedAt               InvoiceProperty = "techInvoicePaidWebhookProcessedAt"
	InvoicePropertyInvoiceFinalizedEventSentAt          InvoiceProperty = "techInvoiceFinalizedSentAt"
	InvoicePropertyPaymentLink                          InvoiceProperty = "paymentLink"
	InvoicePropertyPaymentLinkValidUntil                InvoiceProperty = "paymentLinkValidUntil"
	InvoicePropertyLastRemindInvoiceNotificationSentAt  InvoiceProperty = "lastRemindInvoiceNotificationSentAt"
	InvoicePropertyRemindInvoiceNotificationRequestedAt InvoiceProperty = "techRemindInvoiceNotificationRequestedAt"
	InvoicePropertyPayNotificationRequestedAt           InvoiceProperty = "techPayNotificationRequestedAt"
	InvoicePropertyProviderEmail                        InvoiceProperty = "providerEmail"
	InvoicePropertyProviderCCEmails                     InvoiceProperty = "providerCCEmails"
	InvoicePropertyProviderBCCEmails                    InvoiceProperty = "providerBCCEmails"
	InvoicePropertyPaymentProcessingAt                  InvoiceProperty = "techPaymentProcessingAt"
	InvoicePropertyPaidInvoiceNotificationSentAt        InvoiceProperty = "techPaidInvoiceNotificationSentAt"
	InvoicePropertyVoidInvoiceNotificationRequestAt     InvoiceProperty = "techVoidInvoiceNotificationSentAt"
	InvoicePropertyQuickbooksInvoiceId                  InvoiceProperty = "quickbooksInvoiceId"
	InvoicePropertyQuickbooksJournalEntryId             InvoiceProperty = "quickbooksJournalEntryId"
	InvoicePropertyQuickbooksJournalEntryIdReverse      InvoiceProperty = "quickbooksJournalEntryIdReverse"
	InvoicePropertyQuickbooksPaymentId                  InvoiceProperty = "quickbooksPaymentId"
	InvoicePropertyProviderBankDetailsAvailable         InvoiceProperty = "providerBankDetailsAvailable"
	InvoicePropertyProviderBankAccountName              InvoiceProperty = "providerBankAccountName"
	InvoicePropertyProviderBankAccountNumber            InvoiceProperty = "providerBankAccountNumber"
	InvoicePropertyProviderBankAccountIBAN              InvoiceProperty = "providerBankAccountIBAN"
	InvoicePropertyProviderBankAccountBIC               InvoiceProperty = "providerBankAccountBIC"
	InvoicePropertyProviderBankAccountSortCode          InvoiceProperty = "providerBankAccountSortCode"
	InvoicePropertyProviderBankAccountRoutingNumber     InvoiceProperty = "providerBankAccountRoutingNumber"
	InvoicePropertyProviderBankAccountOtherDetails      InvoiceProperty = "providerBankAccountOtherDetails"
)

type InvoiceEntity struct {
	EventStoreAggregate
	Id                              string
	CreatedAt                       time.Time `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt                       time.Time
	DryRun                          bool          `neo4jDb:"property:dryRun;lookupName:DRY_RUN;supportCaseSensitive:false"`
	Number                          string        `neo4jDb:"property:number;lookupName:NUMBER;supportCaseSensitive:false"`
	Currency                        enum.Currency `neo4jDb:"property:currency;lookupName:CURRENCY;supportCaseSensitive:false"`
	PeriodStartDate                 time.Time     // Date only
	PeriodEndDate                   time.Time     // Date only
	DueDate                         time.Time     `neo4jDb:"property:dueDate;lookupName:DUE_DATE;supportCaseSensitive:false"`       // Date only
	IssuedDate                      time.Time     `neo4jDb:"property:issuedDate;lookupName:ISSUED_DATE;supportCaseSensitive:false"` // Datetime
	Customer                        InvoiceCustomer
	Provider                        InvoiceProvider
	Amount                          float64 `neo4jDb:"property:amount;lookupName:AMOUNT;supportCaseSensitive:false"`
	Vat                             float64 `neo4jDb:"property:vat;lookupName:VAT;supportCaseSensitive:false"`
	TotalAmount                     float64 `neo4jDb:"property:totalAmount;lookupName:TOTAL_AMOUNT;supportCaseSensitive:false"`
	RepositoryFileId                string
	BillingCycleInMonths            int64
	Status                          enum.InvoiceStatus `neo4jDb:"property:status;lookupName:STATUS;supportCaseSensitive:false"`
	Note                            string
	QuickbooksInvoiceId             string
	QuickbooksJournalEntryId        string
	QuickbooksJournalEntryIdReverse string
	QuickbooksPaymentId             string
	PaymentDetails                  PaymentDetails
	OffCycle                        bool
	Postpaid                        bool
	Preview                         bool `neo4jDb:"property:preview;lookupName:PREVIEW;supportCaseSensitive:false"`

	Source    DataSource
	AppSource string

	InvoiceInternalFields InvoiceInternalFields

	DataloaderKey string
}

type InvoiceCustomer struct {
	Name         string
	Email        string
	AddressLine1 string
	AddressLine2 string
	Zip          string
	Locality     string
	Country      string
	Region       string
}

type InvoiceProvider struct {
	LogoRepositoryFileId     string
	Name                     string
	Email                    string
	CC                       []string
	BCC                      []string
	AddressLine1             string
	AddressLine2             string
	Zip                      string
	Locality                 string
	Country                  string
	Region                   string
	BankDetailsAvailable     bool
	BankAccountName          string
	BankAccountNumber        string
	BankAccountIBAN          string
	BankAccountBIC           string
	BankAccountSortCode      string
	BankAccountRoutingNumber string
	BankAccountOtherDetails  string
}

type PaymentDetails struct {
	PaymentLink           string
	PaymentLinkValidUntil *time.Time
}

type InvoiceInternalFields struct {
	InvoiceFinalizedSentAt               *time.Time // used to send the invoice finalized notification to slack and integration app
	InvoiceFinalizedWebhookProcessedAt   *time.Time // used to process webhook for invoice finalized to temporal, if no webhook is configured, property will be set
	InvoicePaidWebhookProcessedAt        *time.Time // used to process webhook for invoice paid to temporal, if no webhook is configured, property will be set
	PaymentLinkRequestedAt               *time.Time
	PayInvoiceNotificationRequestedAt    *time.Time // used for locking in batch to not send the same notification multiple times under an hour
	PayInvoiceNotificationSentAt         *time.Time // used to prevent sending the same notification
	RemindInvoiceNotificationRequestedAt *time.Time
	LastRemindInvoiceNotificationSentAt  *time.Time
	RemindInvoiceNotificationRequestAt   *time.Time
	PaidInvoiceNotificationSentAt        *time.Time
	VoidInvoiceNotificationSentAt        *time.Time
	PaymentProcessingAt                  *time.Time // when invoice was set to payment processing status
}

type InvoiceEntities []InvoiceEntity

func (i InvoiceEntity) IsDue() bool {
	return i.Status == enum.InvoiceStatusDue
}

func (i InvoiceEntity) IsOverdue() bool {
	return i.Status == enum.InvoiceStatusOverdue
}

func (i InvoiceEntity) IsPaid() bool {
	return i.Status == enum.InvoiceStatusPaid
}

func (i InvoiceEntity) IsVoid() bool {
	return i.Status == enum.InvoiceStatusVoid
}

func (i InvoiceEntity) IsPaymentProcessing() bool {
	return i.Status == enum.InvoiceStatusPaymentProcessing
}
