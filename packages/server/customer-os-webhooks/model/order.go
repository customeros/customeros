package model

import "time"

type OrderData struct {
	BaseData
	OrderedByOrganization ReferencedOrganization `json:"orderedByOrganization,omitempty"`
	ConfirmedAt           *time.Time             `json:"confirmedAt"`
	PaidAt                *time.Time             `json:"paidAt"`
	FulfilledAt           *time.Time             `json:"fulfilledAt"`
	CanceledAt            *time.Time             `json:"canceledAt"`
}

func (l *OrderData) Normalize() {
	l.BaseData.Normalize()
	l.SetTimes()
}
