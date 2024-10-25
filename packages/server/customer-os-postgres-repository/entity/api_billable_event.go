package entity

import "time"

// BillableEvent is a custom type representing various events
type BillableEvent string

// Constants for predefined Billable Events
const (
	BillableEventEmailVerifiedCatchAll     BillableEvent = "email_verified_catch_all"
	BillableEventEmailVerifiedNotCatchAll  BillableEvent = "email_verified_not_catch_all"
	BillableEventEnrichPersonEmailFound    BillableEvent = "enrich_person_email_found"
	BillableEventEnrichPersonPhoneFound    BillableEvent = "enrich_person_phone_found"
	BillableEventEnrichOrganizationSuccess BillableEvent = "enrich_organization_success"
	BillableEventIpVerificationSuccess     BillableEvent = "ip_verification_success"
)

// ApiBillableEvent represents a chargeable event in your system
type ApiBillableEvent struct {
	ID            uint64        `gorm:"primary_key;autoIncrement" json:"id"`
	Tenant        string        `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	Event         BillableEvent `gorm:"column:event;type:varchar(255);NOT NULL;index:idx_tenant_event" json:"event"`
	Subtype       string        `gorm:"column:subtype;type:varchar(255)" json:"subtype,omitempty"`
	ExternalID    string        `gorm:"column:external_id;type:varchar(255);index" json:"externalId"`
	ReferenceData string        `gorm:"column:reference_data;type:text" json:"referenceData,omitempty"`
	CreatedAt     time.Time     `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
}

// TableName sets the name of the table for GORM
func (ApiBillableEvent) TableName() string {
	return "api_billable_event"
}
