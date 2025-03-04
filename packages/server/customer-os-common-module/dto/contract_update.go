package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
)

type UpdateContract struct {
	data_fields.ContractSaveFields
	Ltv *float64 `json:"ltv,omitempty"`
}
