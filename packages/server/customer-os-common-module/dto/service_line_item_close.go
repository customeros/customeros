package dto

import "time"

type CloseServiceLineItem struct {
	ServiceLineItemId string    `json:"serviceLineItemId,omitempty"`
	EndedAt           time.Time `json:"endedAt,omitempty"`
}
