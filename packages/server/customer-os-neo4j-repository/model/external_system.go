package model

import "time"

type ExternalSystem struct {
	ExternalSystemId string     `json:"externalSystemId,omitempty"`
	ExternalUrl      string     `json:"externalUrl,omitempty"`
	ExternalId       string     `json:"externalId,omitempty"`
	ExternalIdSecond string     `json:"externalIdSecond,omitempty"`
	ExternalSource   string     `json:"externalSource,omitempty"`
	SyncDate         *time.Time `json:"syncDate,omitempty"`
	Primary          bool       `json:"primary"`
}

func (e ExternalSystem) Available() bool {
	return e.ExternalSystemId != "" && e.ExternalId != ""
}
