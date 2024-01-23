package model

type InvoiceData struct {
	BaseData
	Status string `json:"status,omitempty"`
}

func (i *InvoiceData) Normalize() {
	i.SetTimes()
	i.BaseData.Normalize()
}
