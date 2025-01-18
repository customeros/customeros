package postgres_entity

import "time"

type MailstackBuyRequest struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`

	Domains   string `gorm:"column:domains;type:text" json:"domains"`
	Usernames string `gorm:"column:usernames;type:text" json:"usernames"`

	PaymentIntentId string `gorm:"column:payment_intent_id;type:varchar(255)" json:"paymentIntentId"`

	Status MailstackBuyRequestStatus `gorm:"column:status;type:varchar(50)" json:"status"`
}

func (MailstackBuyRequest) TableName() string {
	return "mailstack_buy_request"
}

type MailstackBuyRequestStatus string

const (
	MailstackBuyRequestStatusPending   MailstackBuyRequestStatus = "PENDING"
	MailstackBuyRequestStatusCompleted MailstackBuyRequestStatus = "COMPLETED"
	MailstackBuyRequestStatusFailed    MailstackBuyRequestStatus = "FAILED"
)

type MailstackBuyRequestDomain struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`

	MailstackBuyRequestId string `gorm:"column:mailstack_buy_request_id;type:uuid;NOT NULL" json:"mailstackBuyRequestId"`
	Domain                string `gorm:"column:domain;type:varchar(255)" json:"domain"`
	RedirectWebsite       string `gorm:"column:redirect_website;type:varchar(255)" json:"redirectWebsite"`

	Status MailstackBuyRequestDomainStatus `gorm:"column:status;type:varchar(50)" json:"status"`
}

func (MailstackBuyRequestDomain) TableName() string {
	return "mailstack_buy_request_domain"
}

type MailstackBuyRequestDomainStatus string

const (
	MailstackBuyRequestDomainStatusPendingProvisioning  MailstackBuyRequestDomainStatus = "PENDING_PROVISIONING"
	MailstackBuyRequestDomainStatusPendingConfiguration MailstackBuyRequestDomainStatus = "PENDING_CONFIGURATION"
	MailstackBuyRequestDomainStatusCompleted            MailstackBuyRequestDomainStatus = "COMPLETED"
	MailstackBuyRequestDomainStatusFailed               MailstackBuyRequestDomainStatus = "FAILED"
)
