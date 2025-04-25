package dto

import "time"

type WebTrackerUpdate struct {
	ID                string
	CNAMEHost         *string
	IsCNAMEConfigured *bool
	IsProxyActive     *bool
	LastEventAt       *time.Time
}
