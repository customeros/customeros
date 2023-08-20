package entity

import "time"

type InteractionSession struct {
	ExternalId   string     `json:"externalId,omitempty"`
	Name         string     `json:"name,omitempty"`
	Channel      string     `json:"channel,omitempty"`
	Type         string     `json:"type,omitempty"`
	CreatedAt    *time.Time `json:"createdAtTime,omitempty"`
	CreatedAtStr string     `json:"createdAt,omitempty"`
	Status       string     `json:"status,omitempty"`
	Identifier   string     `json:"identifier,omitempty"`
}
