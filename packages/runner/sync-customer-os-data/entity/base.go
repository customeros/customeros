package entity

import "time"

type BaseData struct {
	Id             string     `json:"id,omitempty"`
	ExternalId     string     `json:"externalId,omitempty"`
	ExternalSystem string     `json:"externalSystem,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
	SyncId         string     `json:"syncId,omitempty"`
}
