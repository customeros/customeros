package events_platform

import (
	"context"
	invoicepb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
)

type MockInvoiceServiceCallbacks struct {
	VoidInvoice func(context.Context, *invoicepb.VoidInvoiceRequest) (*invoicepb.InvoiceIdResponse, error)
}

var invoiceCallbacks = &MockInvoiceServiceCallbacks{}

func SetInvoiceCallbacks(callbacks *MockInvoiceServiceCallbacks) {
	invoiceCallbacks = callbacks
}

type MockInvoiceService struct {
	invoicepb.UnimplementedInvoiceGrpcServiceServer
}

func (MockInvoiceService) VoidInvoice(context context.Context, proto *invoicepb.VoidInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	if invoiceCallbacks.VoidInvoice == nil {
		panic("invoiceCallbacks.VoidInvoice is not set")
	}
	return invoiceCallbacks.VoidInvoice(context, proto)
}
