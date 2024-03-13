package entity

import "time"

type OrderData struct {
	BaseData
	OrderedByOrganization ReferencedOrganization `json:"orderedByOrganization,omitempty"`
	ConfirmedAt           *time.Time             `json:"confirmedAt,omitempty"`
	PaidAt                *time.Time             `json:"paidAt,omitempty"`
	FulfilledAt           *time.Time             `json:"fulfilledAt,omitempty"`
	CancelledAt           *time.Time             `json:"cancelledAt,omitempty"`
}
