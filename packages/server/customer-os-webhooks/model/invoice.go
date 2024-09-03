package model

type InvoiceData struct {
	BaseData
	Status                string `json:"status,omitempty"`
	PaymentLink           string `json:"paymentLink,omitempty"`
	PaymentLinkValidHours string `json:"paymentLinkValidHours,omitempty"`
}

func (i *InvoiceData) Normalize() {
	i.SetTimes()
	i.BaseData.Normalize()
}
