package events_platform

import (
	"context"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
)

type MockInvoiceServiceCallbacks struct {
	SimulateInvoice func(context.Context, *invoicepb.SimulateInvoiceRequest) (*invoicepb.InvoiceIdResponse, error)
}

var invoiceCallbacks = &MockInvoiceServiceCallbacks{}

func SetInvoiceCallbacks(callbacks *MockInvoiceServiceCallbacks) {
	invoiceCallbacks = callbacks
}

type MockInvoiceService struct {
	invoicepb.UnimplementedInvoiceGrpcServiceServer
}

func (MockInvoiceService) SimulateInvoice(context context.Context, proto *invoicepb.SimulateInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	if invoiceCallbacks.SimulateInvoice == nil {
		panic("invoiceCallbacks.SimulateInvoice is not set")
	}
	return invoiceCallbacks.SimulateInvoice(context, proto)
}
