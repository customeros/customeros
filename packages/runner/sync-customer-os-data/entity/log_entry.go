package entity

import (
	"time"
)

type LogEntryData struct {
	BaseData
	Content              string                 `json:"content,omitempty"`
	ContentType          string                 `json:"contentType,omitempty"`
	StartedAtStr         string                 `json:"startedAt,omitempty"`
	StartedAt            *time.Time             `json:"startedAtTime,omitempty"`
	AuthorUser           ReferencedUser         `json:"authorUser,omitempty"`
	LoggedOrganization   ReferencedOrganization `json:"loggedOrganization,omitempty"`
	LoggedEntityRequired bool                   `json:"loggedEntityRequired,omitempty"`
}
