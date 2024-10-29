package dto

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"

type EventCompleted struct {
	Tenant     string           `json:"tenant"`
	EntityType model.EntityType `json:"entityType"`
	EntityIds  []string         `json:"entityIds"`
	Create     bool             `json:"create"`
	Update     bool             `json:"update"`
	Delete     bool             `json:"delete"`
}
