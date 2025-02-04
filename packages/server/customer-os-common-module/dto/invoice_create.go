package dto

import "time"

type CreateInvoice struct {
	Id                   string    `json:"id"`
	Number               string    `json:"number"`
	Currency             string    `json:"currency"`
	Postpaid             bool      `json:"postpaid"`
	DryRun               bool      `json:"dryRun"`
	Preview              bool      `json:"preview"`
	OffCycle             bool      `json:"offCycle"`
	Note                 string    `json:"note"`
	BillingCycleInMonths int64     `json:"billingCycleInMonths"`
	InvoicePeriodStart   time.Time `json:"invoicePeriodStart"`
	InvoicePeriodEnd     time.Time `json:"invoicePeriodEnd"`
	Source               string    `json:"source"`
	AppSource            string    `json:"appSource"`
	DueDate              time.Time `json:"dueDate"`
	IssuedDate           time.Time `json:"issuedDate"`
}
