package model

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"

type ExternalSystemData struct {
	ExternalSystem     string   `json:"externalSystem"`
	AppSource          string   `json:"appSource"`
	PaymentMethodTypes []string `json:"paymentMethodTypes"`
}

func (e *ExternalSystemData) Normalize() {
	e.PaymentMethodTypes = utils.RemoveEmpties(e.PaymentMethodTypes)
	e.PaymentMethodTypes = utils.RemoveDuplicates(e.PaymentMethodTypes)
}
