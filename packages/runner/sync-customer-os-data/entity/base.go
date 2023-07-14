package entity

import "time"

/*
{
  "skip": false,
  "skipReason": "draft data",
  "id": "1234",
  "externalId": "abcd1234",
  "externalSystem": "HubSpot",
  "createdAt": "2022-02-28T19:52:05Z",
  "updatedAt": "2022-03-01T11:23:45Z",
  "syncId": "sync_1234"
}
*/

type BaseData struct {
	Skip           bool       `json:"skip,omitempty"`
	SkipReason     string     `json:"skipReason,omitempty"`
	Id             string     `json:"id,omitempty"`
	ExternalId     string     `json:"externalId,omitempty"`
	ExternalSystem string     `json:"externalSystem,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
	SyncId         string     `json:"syncId,omitempty"`
}
