package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

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
	CreatedAtStr   string     `json:"createdAt,omitempty"`
	UpdatedAtStr   string     `json:"updatedAt,omitempty"`
	CreatedAt      *time.Time `json:"createdAtTime,omitempty"`
	UpdatedAt      *time.Time `json:"updatedAtTime,omitempty"`
	SyncId         string     `json:"syncId,omitempty"`
}

func (b *BaseData) SetCreatedAt() {
	if b.CreatedAtStr != "" && b.CreatedAt == nil {
		b.CreatedAt, _ = utils.UnmarshalDateTime(b.CreatedAtStr)
	}
}

func (b *BaseData) SetUpdatedAt() {
	if b.UpdatedAtStr != "" && b.UpdatedAt == nil {
		b.UpdatedAt, _ = utils.UnmarshalDateTime(b.UpdatedAtStr)
	}
}

func (b *BaseData) SetTimes() {
	b.SetCreatedAt()
	b.SetUpdatedAt()
	if b.CreatedAt != nil {
		b.CreatedAt = common_utils.TimePtr((*b.CreatedAt).UTC())
	} else {
		b.CreatedAt = common_utils.TimePtr(common_utils.Now())
	}
	if b.UpdatedAt != nil {
		b.UpdatedAt = common_utils.TimePtr((*b.UpdatedAt).UTC())
	} else {
		b.UpdatedAt = common_utils.TimePtr(common_utils.Now())
	}
}
