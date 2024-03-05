package model

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type InteractionSession struct {
	Tenant          string                       `json:"tenant"`
	ID              string                       `json:"id"`
	Source          commonmodel.Source           `json:"source"`
	ExternalSystems []commonmodel.ExternalSystem `json:"externalSystem"`
	CreatedAt       time.Time                    `json:"createdAt"`
	UpdatedAt       time.Time                    `json:"updatedAt"`
	Channel         string                       `json:"channel"`
	ChannelData     string                       `json:"channelData"`
	Identifier      string                       `json:"identifier"`
	Status          string                       `json:"status"`
	Type            string                       `json:"type"`
	Name            string                       `json:"name"`
}
