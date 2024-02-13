package model

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strings"
	"time"
)

type BaseData struct {
	Skip                 bool       `json:"skip,omitempty"`
	SkipReason           string     `json:"skipReason,omitempty"`
	Id                   string     `json:"id,omitempty"`
	ExternalId           string     `json:"externalId,omitempty"`
	ExternalIdSecond     string     `json:"externalIdSecond,omitempty"`
	ExternalSystem       string     `json:"externalSystem,omitempty"`
	ExternalUrl          string     `json:"externalUrl,omitempty"`
	ExternalSourceEntity string     `json:"externalSourceEntity,omitempty"`
	CreatedAtStr         string     `json:"createdAt,omitempty"`
	UpdatedAtStr         string     `json:"updatedAt,omitempty"`
	CreatedAt            *time.Time `json:"createdAtTime,omitempty"`
	UpdatedAt            *time.Time `json:"updatedAtTime,omitempty"`
	SyncId               string     `json:"syncId,omitempty"`
	AppSource            string     `json:"appSource,omitempty"`
	Source               string     `json:"source,omitempty"` // used if ExternalSystem is empty
	UpdateOnly           bool       `json:"updateOnly"`
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
		b.CreatedAt = utils.TimePtr((*b.CreatedAt).UTC())
	} else {
		b.CreatedAt = utils.TimePtr(utils.Now())
	}
	if b.UpdatedAt != nil {
		b.UpdatedAt = utils.TimePtr((*b.UpdatedAt).UTC())
	} else {
		b.UpdatedAt = utils.NowPtr()
	}
}

func (b *BaseData) Normalize() {
	if b.ExternalSystem != "" {
		b.ExternalSystem = strings.ToLower(strings.TrimSpace(b.ExternalSystem))
	}
	if b.ExternalId != "" {
		b.ExternalId = strings.TrimSpace(b.ExternalId)
	}
}
