package data_fields

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type ContractSaveFields struct {
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	AppSource       *string    `json:"appSource,omitempty"`
	Source          *string    `json:"source,omitempty"`
	OrganizationId  *string    `json:"organizationId,omitempty"`
	CreatedByUserId *string    `json:"createdByUserId,omitempty"`
}

func (c ContractSaveFields) GetCreatedByUserId() string {
	return utils.IfNotNilString(c.CreatedByUserId)
}

func (c ContractSaveFields) GetOrganizationId() string {
	return utils.IfNotNilString(c.OrganizationId)
}
