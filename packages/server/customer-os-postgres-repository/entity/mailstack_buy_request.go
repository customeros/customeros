package entity

import "time"

type MailstackBuyRequest struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`

	Domains   string `gorm:"column:domains;type:text" json:"domains"`
	Usernames string `gorm:"column:usernames;type:text" json:"usernames"`

	PaymentIntentId           string `gorm:"column:payment_intent_id;type:varchar(255)" json:"paymentIntentId"`
	PaymentIntentClientSecret string `gorm:"column:payment_intent_client_secret;type:varchar(255)" json:"paymentIntentClientSecret"`

	Status MailstackBuyRequestStatus `gorm:"column:status;type:varchar(50)" json:"status"`
}

func (MailstackBuyRequest) TableName() string {
	return "mailstack_buy_request"
}

type MailstackBuyRequestStatus string

const (
	MailstackBuyRequestStatusAwaitingPayment MailstackBuyRequestStatus = "AWAITING_PAYMENT"
	MailstackBuyRequestStatusPending         MailstackBuyRequestStatus = "PENDING"
	MailstackBuyRequestStatusCompleted       MailstackBuyRequestStatus = "COMPLETED"
	MailstackBuyRequestStatusFailed          MailstackBuyRequestStatus = "FAILED"
)

type MailstackBuyRequestDomain struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`

	MailstackBuyRequestId string `gorm:"column:mailstack_buy_request_id;type:uuid;NOT NULL" json:"mailstackBuyRequestId"`
	Domain                string `gorm:"column:domain;type:varchar(255)" json:"domain"`

	Status MailstackBuyRequestDomainStatus `gorm:"column:status;type:varchar(50)" json:"status"`
}

func (MailstackBuyRequestDomain) TableName() string {
	return "mailstack_buy_request_domain"
}

type MailstackBuyRequestDomainStatus string

const (
	MailstackBuyRequestDomainStatusAwaitingPayment MailstackBuyRequestDomainStatus = "AWAITING_PAYMENT"
	MailstackBuyRequestDomainStatusPending         MailstackBuyRequestDomainStatus = "PENDING"
	MailstackBuyRequestDomainStatusCompleted       MailstackBuyRequestDomainStatus = "COMPLETED"
	MailstackBuyRequestDomainStatusFailed          MailstackBuyRequestDomainStatus = "FAILED"
)

type MailstackBuyRequestMailbox struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`

	MailstackBuyRequestId string `gorm:"column:mailstack_buy_request_id;type:uuid;NOT NULL" json:"mailstackBuyRequestId"`
	Domain                string `gorm:"column:domain;type:varchar(255)" json:"domain"`
	Username              string `gorm:"column:username;type:varchar(255)" json:"username"`
	Mailbox               string `gorm:"column:mailbox;type:varchar(255)" json:"mailbox"`

	Status MailstackBuyRequestMailboxStatus `gorm:"column:status;type:varchar(50)" json:"status"`
}

func (MailstackBuyRequestMailbox) TableName() string {
	return "mailstack_buy_request_mailbox"
}

type MailstackBuyRequestMailboxStatus string

const (
	MailstackBuyRequestDomainMailboxAwaitingPayment MailstackBuyRequestMailboxStatus = "AWAITING_PAYMENT"
	MailstackBuyRequestDomainMailboxPending         MailstackBuyRequestMailboxStatus = "PENDING"
	MailstackBuyRequestDomainMailboxCompleted       MailstackBuyRequestMailboxStatus = "COMPLETED"
	MailstackBuyRequestDomainMailboxFailed          MailstackBuyRequestMailboxStatus = "FAILED"
)
