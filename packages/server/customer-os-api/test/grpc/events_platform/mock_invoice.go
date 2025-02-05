package events_platform

import (
	invoicepb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
)

type MockInvoiceServiceCallbacks struct {
}

var invoiceCallbacks = &MockInvoiceServiceCallbacks{}

func SetInvoiceCallbacks(callbacks *MockInvoiceServiceCallbacks) {
	invoiceCallbacks = callbacks
}

type MockInvoiceService struct {
	invoicepb.UnimplementedInvoiceGrpcServiceServer
}
