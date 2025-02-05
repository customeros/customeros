package events_platform

import ()

type MockInvoiceServiceCallbacks struct {
}

var invoiceCallbacks = &MockInvoiceServiceCallbacks{}

func SetInvoiceCallbacks(callbacks *MockInvoiceServiceCallbacks) {
	invoiceCallbacks = callbacks
}

type MockInvoiceService struct {
}
